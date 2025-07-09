package ftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const TIMEOUT_LOAD_SEC = 17

type FTPSender struct {
	Host       string
	User       string
	Pass       string
	LocalPath  string
	RemotePath string
	Port       int
	Enable     bool
	sqlRepo    persistent.SQLRepo
	cache      temporary.CacheRepo
}

func (f *FTPSender) Loader(ctx context.Context) {
	t := time.NewTicker(time.Second * TIMEOUT_LOAD_SEC)
	for {
		select {
		case <-t.C:
			var err error
			var link *persistent.LinkModel
			if link, err = f.sqlRepo.SelectOne(entity.SENDING); err != nil {
				fmt.Printf("select file to send: %s\n", err.Error())
				break
			}
			if link == nil {
				fmt.Println("file to send not found")
				break
			}

			if err := f.presend(link); err != nil {
				fmt.Printf("send file to sftp: %s\n", err.Error())
				break
			}

			fmt.Printf("file %s sended\n", *link.Filename)
		case <-ctx.Done():
			fmt.Println("context failed")
			return
		}
	}
}

func (f *FTPSender) presend(link *persistent.LinkModel) error {
	f.cache.SetStatus(&temporary.TaskRequest{
		FileName:   *link.Filename,
		Link:       link.Link,
		MoveTo:     f.RemotePath,
		MaxQuality: link.TargetQuantity,
		Procentage: 100,
		Status:     entity.SENDING,
	})

	if f.Enable {
		if err := f.send(*link.Filename, link.Link, link.TargetQuantity); err != nil {
			fmt.Printf("[send file by ftp: %s\\n", err.Error())
			return fmt.Errorf("send file: %s", err.Error())
		}
	}

	f.sqlRepo.UpdateStatus(link.Link, entity.DONE)
	f.cache.LinkDone(link.Link)

	return nil
}

func (f *FTPSender) send(filename, link string, targetQuantity int) error {
	config := &ssh.ClientConfig{
		User: f.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(f.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use a proper HostKeyCallback in production
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", f.Host, f.Port), config)
	if err != nil {
		return fmt.Errorf("ftp send (dial): %w", err)
	}
	defer client.Close()

	sc, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("ftp send (newClient): %w", err)
	}
	defer sc.Close()

	srcFile, err := os.Open(fmt.Sprintf("%s/%s", f.LocalPath, filename))
	if err != nil {
		return fmt.Errorf("ftp send (open): %w", err)
	}
	defer srcFile.Close()

	remoteDir := filepath.Dir(f.RemotePath)
	_ = sc.MkdirAll(remoteDir)

	dstFile, err := sc.Create(fmt.Sprintf("%s/%s", f.RemotePath, filename))
	if err != nil {
		return fmt.Errorf("ftp send (create remote): %w", err)
	}
	defer dstFile.Close()

	st, _ := srcFile.Stat()
	f.cache.SetStatus(&temporary.TaskRequest{
		FileName:     filename,
		Link:         link,
		MoveTo:       f.RemotePath,
		MaxQuality:   targetQuantity,
		Procentage:   100,
		Status:       entity.SENDING,
		DownloadSize: float64(float64(st.Size()/1024) / 1024),
		CurrentSize:  float64(float64(st.Size()/1024) / 1024),
		Message:      "sending",
	})

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("ftp send (copy): %w", err)
	}

	if err = os.Remove(fmt.Sprintf("%s/%s", f.LocalPath, filename)); err != nil {
		return fmt.Errorf("file remove: %w", err)
	}

	return nil
}
