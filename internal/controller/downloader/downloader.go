package downloader

import (
	"context"
	"dwld-downloader/internal/entity"
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/lrstanley/go-ytdlp"
)

const (
	TIMEOUT_WORKERS = 10
	MIN_QUALITY     = 720
	MAX_QUALITY     = 10000
)

type Task struct {
	Link    string
	Quality int
}
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
	PercentToNext int
	Stages        map[int]Stage
	sqlRepo       persistent.SQLRepo
	cache         temporary.CacheRepo
	workersPool   chan struct{}
	totalStages   int
}

type DownloaderConf struct {
	WorkDir       string
	Threads       int
	PercentToNext int
	Stages        []Stage
	SqlRepo       persistent.SQLRepo
	Cache         temporary.CacheRepo
}

func InstallYTDLP() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("YTDLP failed auto install: %v\nLaunch without\n", r)
		}
	}()
	ytdlp.MustInstall(context.Background(), nil)
}

func NewDownloader(conf DownloaderConf) Downloader {
	InstallYTDLP()

	stages := make(map[int]Stage)
	for _, v := range conf.Stages {
		stages[v.Positions] = v
	}

	return &DownloaderSource{
		WorkDir:       conf.WorkDir,
		PercentToNext: conf.PercentToNext,
		Stages:        stages,
		sqlRepo:       conf.SqlRepo,
		cache:         conf.Cache,
		workersPool:   make(chan struct{}, conf.Threads),
		totalStages:   len(stages),
	}
}

func (d *DownloaderSource) Downloader(task *Task) error {
	d.sqlRepo.UpdateStatus(task.Link, entity.WORK)

	qualtiy := task.Quality
	if qualtiy > MAX_QUALITY {
		qualtiy = MAX_QUALITY
	}
	if qualtiy < MIN_QUALITY {
		qualtiy = MIN_QUALITY
	}

	filename := ""

	var toNext sync.Once

	size := float64(0)
	totalSize := float64(0)

	dl := ytdlp.New().
		SetWorkDir(d.WorkDir).
		FormatSort("res,ext:mp4:m4a").
		RecodeVideo("mp4").
		Output("%(title)s.%(ext)s").
		NoRestrictFilenames().
		Fixup(ytdlp.FixupForce).
		Retries("20").
		NoWriteSubs().
		IgnoreErrors().
		IgnoreNoFormatsError().
		NoAbortOnError().
		RmCacheDir().
		Format(fmt.Sprintf("bv*[height<=%d]+ba", qualtiy)).ProgressFunc(time.Duration(time.Millisecond*750), func(update ytdlp.ProgressUpdate) {
		size = (float64(update.DownloadedBytes) / 1024) / 1024 // К мегабайтам
		totalSize = (float64(update.TotalBytes) / 1024) / 1024 // К мегабайтам
		fmt.Println(update.Status, update.PercentString(), fmt.Sprintf("[%.2f/%.2f]mb", size, totalSize), update.Filename)

		status := string(update.Status)
		if strings.Contains(status, "finished") {
			status = "converting"
		}

		//update.Filename
		filename = *update.Info.Filename

		d.cache.SetStatus(&temporary.TaskRequest{
			FileName:     filename,
			Link:         task.Link,
			MoveTo:       d.WorkDir,
			MaxQuality:   qualtiy,
			Procentage:   update.Percent(),
			Status:       entity.WORK,
			DownloadSize: totalSize,
			CurrentSize:  size,
			Message:      status,
		})
		if update.Percent() > float64(d.PercentToNext) {
			toNext.Do(func() {
				<-d.workersPool
			})
		}

	})

	if err := func() error {
		var err_resp error
		for i := range d.totalStages {
			stg := d.Stages[(i + 1)]

			dl.UnsetCookiesFromBrowser()
			if stg.IsCookie {
				dl.CookiesFromBrowser("chrome")
			}

			dl.UnsetEmbedChapters()
			if stg.IsEmbededCharters {
				dl.EmbedChapters()
			}

			dl.UnsetMarkWatched()
			if stg.IsMarkWatched {
				dl.MarkWatched()
			}

			dl.UnsetExtractorArgs()
			if stg.Extractors != "" {
				dl.ExtractorArgs(stg.Extractors)
			}

			for retry := range stg.AttemptBeforeNext {
				_ = retry
				_, err := dl.Run(context.TODO(), task.Link)
				if err != nil {
					d.sqlRepo.Update(&persistent.LinkModelRequest{
						Link:           task.Link,
						Filename:       pointer.To(filename),
						WorkStatus:     entity.WORK,
						Message:        pointer.To(err.Error()),
						TargetQuantity: task.Quality,
					})
					d.cache.SetStatus(&temporary.TaskRequest{
						FileName:     filename,
						Link:         task.Link,
						MoveTo:       d.WorkDir,
						MaxQuality:   qualtiy,
						Procentage:   0,
						Status:       entity.WORK,
						DownloadSize: totalSize,
						CurrentSize:  size,
						Message:      err.Error(),
					})
					err_resp = err
					fmt.Printf("download failed: %s", err.Error())
					continue
				}

				// Скачивание и конвертация прошли успешно
				return nil
			}
		}

		return err_resp
	}(); err != nil {
		return fmt.Errorf("downloader: %w", err)
	}

	return nil
}

func (d *DownloaderSource) Processor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context canceled")
			return
		case d.workersPool <- struct{}{}:
			go func() {
				task, err := d.GetLink()
				if err != nil {
					fmt.Printf("downloader(select link): %s\n", err.Error())
					<-d.workersPool
					return
				}
				if task.Link != "" {
					err := d.Downloader(task)
					if err != nil {
						<-d.workersPool
						fmt.Printf("downloader: %s\n", err.Error())
						// Помещаем обратно в пул
						d.sqlRepo.UpdateStatus(task.Link, entity.NEW)
						d.cache.LinkDone(task.Link)

						return
					}

					d.sqlRepo.UpdateStatus(task.Link, entity.SENDING)
				} else {
					time.Sleep(time.Second * TIMEOUT_WORKERS)
				}
			}()
		default:
			time.Sleep(time.Second * TIMEOUT_WORKERS)
		}
	}
}

func (d *DownloaderSource) GetLink() (*Task, error) {
	link, err := d.sqlRepo.SelectOne(entity.NEW)
	if err != nil {
		return nil, fmt.Errorf("Downloader GetLink (select): %w", err)
	}
	if link == nil {
		link = &persistent.LinkModel{
			Link:           "",
			TargetQuantity: MIN_QUALITY,
		}
	}
	return &Task{
		Link:    link.Link,
		Quality: link.TargetQuantity,
	}, nil
}
