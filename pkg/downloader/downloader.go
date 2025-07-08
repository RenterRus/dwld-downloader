package downloader

import (
	"context"
	"dwld-downloader/internal/controller/downloader"
	"dwld-downloader/internal/entity"
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
)

type DownloaderConf struct {
	WorkDir       string
	Threads       int
	PercentToNext int
	Stages        []entity.Stage
	SqlRepo       persistent.SQLRepo
	Cache         temporary.CacheRepo
}

type Downloader struct {
	downloader downloader.Downloader
	notify     chan struct{}
}

func NewDownloader(conf DownloaderConf) *Downloader {
	return &Downloader{
		downloader: downloader.NewDownloader(downloader.DownloaderConf{
			WorkDir:       conf.WorkDir,
			Threads:       conf.Threads,
			PercentToNext: conf.PercentToNext,
			Stages:        conf.Stages,
			SqlRepo:       conf.SqlRepo,
			Cache:         conf.Cache,
		}),
		notify: make(chan struct{}, 1),
	}
}

func (d *Downloader) Start() {
	ctx, cncl := context.WithCancel(context.Background())
	go func() {
		d.downloader.Processor(ctx)
	}()

	<-d.notify
	cncl()
}

func (d *Downloader) Stop() {
	d.notify <- struct{}{}
}
