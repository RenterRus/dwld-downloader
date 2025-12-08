package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/RenterRus/dwld-downloader/internal/controller/grpc"
	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/register"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	"github.com/RenterRus/dwld-downloader/internal/usecase/download"
	"github.com/RenterRus/dwld-downloader/pkg/cache"
	dwnld "github.com/RenterRus/dwld-downloader/pkg/downloader"
	"github.com/RenterRus/dwld-downloader/pkg/grpcserver"
	"github.com/RenterRus/dwld-downloader/pkg/loader"
	"github.com/RenterRus/dwld-downloader/pkg/sqldb"
	"github.com/samber/lo"
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

	dbconn := sqldb.NewDB(conf.PathToDB, conf.NameDB)
	db := persistent.NewSQLRepo(dbconn, conf.Downloader.WorkPath)
	cc := cache.NewCache(conf.Cache.Host, conf.Cache.Port)
	cache := temporary.NewMemCache(cc)

	dwld := dwnld.NewDownloader(dwnld.DownloaderConf{
		WorkDir:       conf.Downloader.WorkPath,
		Threads:       conf.Downloader.Threads,
		PercentToNext: conf.Downloader.PercentToNext,
		Stages: lo.Map(conf.Downloader.Stages, func(stage Stage, _ int) entity.Stage {
			return entity.Stage{
				Positions:         stage.Positions,
				AttemptBeforeNext: stage.AttemptBeforeNext,
				Threads:           stage.Threads,
				IsCookie:          stage.IsCookie,
				IsEmbededCharters: stage.IsEmbededCharters,
				IsMarkWatched:     stage.IsEmbededCharters,
				Extractors:        stage.Extractors,
				IsFormat:          stage.IsFormat,
			}
		}),
		SqlRepo:   db,
		Cache:     cache,
		EagleMode: true,
	})

	ftpSender := loader.NewLoader(loader.Server{
		Host:    conf.FTP.Host,
		Port:    conf.FTP.Port,
		SqlRepo: db,
		Cache:   cache,
	})

	downloadUsecases := download.NewDownload(
		db,
		cache,
		ftpSender.Sender(),
	)

	// FTPSender
	ctx, cncl := context.WithCancel(context.Background())
	go cache.Revisor(ctx)
	go dwld.Start()
	go ftpSender.Start()

	go func() {
		register.Register(register.RegisterConfig{
			To: register.Server{
				Host: conf.Register.To.Host,
				Port: conf.Register.To.Port,
			},
			From: register.Server{
				Host: conf.GRPC.Host,
				Port: conf.GRPC.Port,
			},
			Assign: conf.Register.Assign,
			Name:   conf.Register.Name,
		})
	}()
	// gRPC Server
	grpcServer := grpcserver.New(grpcserver.Port(conf.GRPC.Host, strconv.Itoa(conf.GRPC.Port)))
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

	cncl()
	cc.Close()
	dwld.Stop()
	ftpSender.Stop()
	dbconn.Close()
	err = grpcServer.Shutdown()

	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - grpcServer.Shutdown: %w", err))
	}

	return nil
}
