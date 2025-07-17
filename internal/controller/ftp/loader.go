package ftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	TIMEOUT_LOAD_SEC = 17
)

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
			if link, err = f.sqlRepo.SelectOne(entity.TO_SEND); err != nil {
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
			fmt.Printf("send file by ftp: %s\\n", err.Error())
			fmt.Printf("attempt with hard mp4")

			filename := strings.Builder{}
			names := strings.Split(*link.Filename, ".")[:len(strings.Split(*link.Filename, "."))-1]
			for i, v := range names {
				if i > 0 {
					filename.WriteString(".")
				}
				filename.WriteString(v)
			}
			filename.WriteString(".mp4")

			if err := f.send(filename.String(), link.Link, link.TargetQuantity); err != nil {
				fmt.Printf("send file by ftp (with mp4): %s\\n", err.Error())

				return fmt.Errorf("send file: %s", err.Error())
			}
		}
	}

	f.sqlRepo.UpdateStatus(link.Link, entity.DONE)
	f.cache.LinkDone(link.Link)

	return nil
}

func (f *FTPSender) send(filename, link string, targetQuantity int) error {
	fmt.Println("Prepare to send (ftp)")
	config := &ssh.ClientConfig{
		User: f.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(f.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use a proper HostKeyCallback in production
	}

	fmt.Println("Prepare connect to ftp")
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", f.Host, f.Port), config)
	if err != nil {
		return fmt.Errorf("ftp send (dial): %w", err)
	}
	defer client.Close()

	fmt.Println("Prepare new client (ftp)")
	sc, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("ftp send (newClient): %w", err)
	}
	defer sc.Close()

	fmt.Println("Open local file (ftp):", fmt.Sprintf("%s/%s", f.LocalPath, filename))
	srcFile, err := os.Open(fmt.Sprintf("%s/%s", f.LocalPath, filename))
	if err != nil {
		return fmt.Errorf("ftp send (open): %w", err)
	}
	defer srcFile.Close()

	fmt.Println("Remote dir (ftp)")
	remoteDir := filepath.Dir(f.RemotePath)
	_ = sc.MkdirAll(remoteDir)

	fmt.Println("Create remote dir (ftp):", fmt.Sprintf("%s/%s", f.RemotePath, filename))
	dstFile, err := sc.Create(fmt.Sprintf("%s/%s", f.RemotePath, filename))
	if err != nil {
		return fmt.Errorf("ftp send (create remote): %w", err)
	}
	defer dstFile.Close()

	st, _ := srcFile.Stat()

	notify := make(chan struct{}, 1)
	go func() {
		t := time.NewTicker(TIMEOUT_LOAD_SEC * time.Second)
		f.cache.SetStatus(&temporary.TaskRequest{
			FileName:     filename,
			Link:         link,
			MoveTo:       f.RemotePath,
			MaxQuality:   targetQuantity,
			Procentage:   100,
			Status:       entity.SENDING,
			DownloadSize: float64(float64(st.Size()/1024) / 1024),
			CurrentSize:  0,
			Message:      "sending",
		})
		for {
			select {
			case <-t.C:
				rmFile, err := sc.OpenFile(fmt.Sprintf("%s/%s", f.RemotePath, filename), os.O_RDONLY)
				if err != nil {
					fmt.Printf("ftp send (OpenFile): %s", err.Error())
				}

				rmStat, err := rmFile.Stat()
				if err != nil {
					fmt.Printf("ftp send (Stat): %s", err.Error())
					fmt.Printf("Sending via ftp [%.2fmb][%s] %s\n", float64(float64(st.Size()/1024)/1024), time.Now().Format(time.DateTime), filename)
					f.cache.SetStatus(&temporary.TaskRequest{
						FileName:     filename,
						Link:         link,
						MoveTo:       f.RemotePath,
						MaxQuality:   targetQuantity,
						Procentage:   100,
						Status:       entity.SENDING,
						DownloadSize: float64(float64(st.Size()/1024) / 1024),
						CurrentSize:  float64(float64(st.Size()/1024) / 1024),
						Message:      "sending (without detail stat)",
					})
				} else {
					curSize := float64(float64(rmStat.Size()/1024) / 1024)
					totalSize := float64(float64(st.Size()/1024) / 1024)
					fmt.Printf("Sending via ftp [%.2f%%][%.2f/%.2fmb][%s] %s\n", (curSize/totalSize)*100.0, curSize, totalSize, time.Now().Format(time.DateTime), filename)
					f.cache.SetStatus(&temporary.TaskRequest{
						FileName:     filename,
						Link:         link,
						MoveTo:       f.RemotePath,
						MaxQuality:   targetQuantity,
						Procentage:   (curSize / totalSize) * 100.0,
						Status:       entity.SENDING,
						DownloadSize: totalSize,
						CurrentSize:  curSize,
						Message:      "sending",
					})
				}
				rmFile.Close()

			case <-notify:
				f.cache.LinkDone(link)
				return
			}
		}
	}()

	fmt.Println("Copy file to remote (ftp):", fmt.Sprintf("%s/%s", f.RemotePath, filename))

	f.sqlRepo.UpdateStatus(link, entity.SENDING)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("ftp send (copy): %w", err)
	}

	fmt.Println("Remove local file (ftp):", fmt.Sprintf("%s/%s", f.LocalPath, filename))
	if err = os.Remove(fmt.Sprintf("%s/%s", f.LocalPath, filename)); err != nil {
		return fmt.Errorf("file remove: %w", err)
	}

	notify <- struct{}{}
	notify <- struct{}{} // lazy wait cancel gorutine
	return nil
}
