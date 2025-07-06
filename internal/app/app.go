package app

import (
	"dwld-downloader/internal/controller/grpc"
	"dwld-downloader/internal/controller/http"
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
	"dwld-downloader/internal/usecase/download"
	"dwld-downloader/pkg/cache"
	"dwld-downloader/pkg/ftp"
	"dwld-downloader/pkg/grpcserver"
	"dwld-downloader/pkg/httpserver"
	"dwld-downloader/pkg/sqldb"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func NewApp(configPath string) error {
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
		httpserver.NewHttpServer(&httpserver.Server{
			Host:   conf.HTTP.Host,
			Port:   conf.HTTP.Port,
			Enable: conf.HTTP.Enable,
			Mux:    http.NewRoute(),
		})
	}()

	db := persistent.NewSQLRepo(sqldb.NewDB(conf.PathToDB, conf.NameDB), conf.Downloader.WorkPath)
	cc := cache.NewCache(conf.Cache.Host, conf.Cache.Port)
	cache := temporary.NewMemCache(cc)

	downloadUsecases := download.NewDownload(
		db,
		cache,
	)

	// FTPSender
	ftpSender := ftp.NewSender(ftp.SenderConf{
		Host:       conf.FTP.Addr.Host,
		User:       conf.FTP.User,
		Pass:       conf.FTP.Pass,
		LocalPath:  conf.Downloader.WorkPath,
		RemotePath: conf.FTP.RemoteDirectory,
		Port:       conf.FTP.Addr.Port,
		SqlRepo:    db,
		Cache:      cache,
	})
	ftpSender.Start()

	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(strconv.Itoa(conf.GRPC.Port)))
	grpc.NewRouter(grpcServer.App, downloadUsecases)
	grpcServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal: %s\n", s.String())
	case err = <-grpcServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Notify: %w", err))
	}

	cc.Close()
	ftpSender.Stop()
	err = grpcServer.Shutdown()

	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
