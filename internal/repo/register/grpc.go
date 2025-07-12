package register

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/RenterRus/dwld-bot/docs/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Host string
	Port int
}

type RegisterConfig struct {
	To     Server
	From   Server
	Assign string
	Name   string
}

func Register(conf RegisterConfig) {
	cc, err := grpc.NewClient(fmt.Sprintf("%s:%d", conf.To.Host, conf.To.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("RegisterDownloader(Client): %s\n", err.Error())
	}
	_, err = v1.NewBotClient(cc).RegisterDownloader(context.Background(), &v1.RegisterDownloaderRequest{
		ServerName:       conf.Name,
		ServerHost:       conf.From.Host,
		ServerPort:       int32(conf.From.Port),
		AllowedRootLinks: strings.Split(conf.Assign, ","),
	})
	if err != nil {
		fmt.Printf("RegisterDownloader(RegisterDownloader): %s\n", err.Error())

	}
}
