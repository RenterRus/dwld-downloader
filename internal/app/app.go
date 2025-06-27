package app

import (
	"dwld-downloader/internal/controller/grpc"
	"dwld-downloader/internal/controller/http"
	"dwld-downloader/internal/usecase/download"
	"dwld-downloader/pkg/grpcserver"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func NewApp(configPath string) error {
	l := log.Logger{}
	lastSlash := 0
	for i, v := range configPath {
		if v == '/' {
			lastSlash = i
		}
	}

	conf, err := ReadConfig(configPath[:lastSlash], configPath[lastSlash+1:])
	if err != nil {
		return fmt.Errorf("ReadConfig: %w", err)
	}

	go func() {
		http.NewHttpServer(&http.Server{
			Host:   conf.HTTP.Host,
			Port:   conf.HTTP.Port,
			Enable: conf.HTTP.Enable,
		})
	}()

	downloadUsecases := download.NewDownload()
	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases, l)

	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		l.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	err = grpcServer.Shutdown()
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
