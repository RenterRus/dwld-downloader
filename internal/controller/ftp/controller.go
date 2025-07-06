package ftp

import (
	"fmt"
	"os"

	ftpapi "github.com/moov-io/go-ftp"
)

type Sender interface {
	Send(filename string) error
}

type FTPSender struct {
	Host       string
	User       string
	Pass       string
	LocalPath  string
	RemotePath string
	Port       int
}

func (f FTPSender) Send(filename string) error {
	clientConfig := ftpapi.ClientConfig{
		Hostname: fmt.Sprintf("%s:%d", f.Host, f.Port),
		Username: f.User,
		Password: f.Pass,
	}

	fmt.Println("1")
	client, err := ftpapi.NewClient(clientConfig)
	if err != nil {
		return fmt.Errorf("new client: %w", err)
	}

	fmt.Println("2")
	// Check if the FTP client is reachable
	if err := client.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	fmt.Println("3")
	// Open the file to be uploaded
	fileData, err := os.Open(fmt.Sprintf("%s/%s", f.LocalPath, filename))
	if err != nil {
		return fmt.Errorf("open file error: %w", err)
	}

	fmt.Println("4")
	// Upload the file to the destination path
	err = client.UploadFile(fmt.Sprintf("%s/%s", f.RemotePath, filename), fileData)
	if err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	fmt.Println("5")
	// Remove from local storage after succ send
	err = os.Remove(fmt.Sprintf("%s/%s", f.LocalPath, filename))
	if err != nil {
		return fmt.Errorf("FTP Send (remove from local): %w", err)
	}

	return nil
}
