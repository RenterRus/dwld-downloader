package downloader

import (
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
)

type Stage struct {
	Positions         int
	AttemptBeforeNext int
	Threads           int
	IsCookie          bool
	IsEmbededCharters bool
	IsMarkWatched     bool
	Extractors        string
}

type DownloaderSource struct {
	WorkDir       string
	Threads       int
	PercentToNext int
	Stages        []Stage
	sqlRepo       persistent.SQLRepo
	cache         temporary.CacheRepo
}

type DownloaderConf struct {
	WorkDir       string
	Threads       int
	PercentToNext int
	Stages        []Stage
	SqlRepo       persistent.SQLRepo
	Cache         temporary.CacheRepo
}

func NewDownloader(conf DownloaderConf) Downloader {
	return &DownloaderSource{
		WorkDir:       conf.WorkDir,
		Threads:       conf.Threads,
		PercentToNext: conf.PercentToNext,
		Stages:        conf.Stages,
		sqlRepo:       conf.SqlRepo,
		cache:         conf.Cache,
	}
}

func (d *DownloaderSource) Downloader(link string, quality int32) {

}

func (d *DownloaderSource) Processor() {

}

func (d *DownloaderSource) GetLink() (string, error) {

}
