package loader

import (
	"context"
	"fmt"

	"github.com/RenterRus/dwld-downloader/internal/controller/ftp"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	v1 "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Host    string
	Port    int
	SqlRepo persistent.SQLRepo
	Cache   temporary.CacheRepo
}

type Loader struct {
	sender ftp.Sender
	notify chan struct{}
}

func NewLoader(conf Server) Loader {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.Host, conf.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("RegisterDownloader(Client): %s\n", err.Error())
	}

	return Loader{
		sender: ftp.NewFTPLoader(v1.NewSenderClient(cc), conf.SqlRepo, conf.Cache),
		notify: make(chan struct{}, 1),
	}
}

func (l *Loader) Sender() ftp.Sender {
	return l.sender
}

func (l *Loader) Start() {
	ctx := context.Background()
	go l.sender.Sender(ctx)
	<-l.notify
	ctx.Done()
}

func (l *Loader) Stop() {
	l.notify <- struct{}{}
}
