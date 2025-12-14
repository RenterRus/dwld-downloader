package downloader

import (
	"context"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"

	"github.com/lrstanley/go-ytdlp"
)

const (
	TIMEOUT_WORKERS = 7
	MIN_QUALITY     = 720
	MAX_QUALITY     = 10000
)

type Task struct {
	Link    string
	Quality int
}

type DownloaderSource struct {
	WorkDir       string
	PercentToNext int
	Stages        map[int]entity.Stage
	sqlRepo       persistent.SQLRepo
	cache         temporary.CacheRepo
	workersPool   chan struct{}
	totalStages   int
	eagleMode     bool
}

type DownloaderConf struct {
	WorkDir       string
	Threads       int
	PercentToNext int
	Stages        []entity.Stage
	SqlRepo       persistent.SQLRepo
	Cache         temporary.CacheRepo
	EagleMode     bool
}

func MustInstallTools() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Validate tools install: %v\n", r)
		}
	}()

	res := ytdlp.MustInstallAll(context.Background())
	fmt.Println("==================")
	for i, v := range res {
		if i > 0 {
			fmt.Println("-------------------------")
		}
		fmt.Printf("Executable: %s\nVersion:    %s\nFromCache:  %v\nDownloaded: %v\n", v.Executable, v.Version, v.FromCache, v.Downloaded)
	}
	fmt.Println("==================")

	ytdlp.MustInstall(context.Background(), nil)
}

func NewDownloader(conf DownloaderConf) Downloader {
	MustInstallTools()

	stages := make(map[int]entity.Stage)
	for _, v := range conf.Stages {
		stages[v.Positions] = v
	}

	return &DownloaderSource{
		WorkDir: conf.WorkDir,
		PercentToNext: func() int {
			if conf.PercentToNext < 1 {
				return 1
			}
			if conf.PercentToNext > 99 {
				return 99
			}

			return conf.PercentToNext
		}(),
		Stages:      stages,
		sqlRepo:     conf.SqlRepo,
		cache:       conf.Cache,
		workersPool: make(chan struct{}, conf.Threads),
		totalStages: len(stages),
	}
}

type statInfo struct {
	task        *Task
	msg         string
	filename    string
	totalSize   float64
	currectSize float64
	procentage  float64
	status      entity.Status
}

func (d *DownloaderSource) statUpdate(stat statInfo) {
	if err := d.sqlRepo.Update(&persistent.LinkModelRequest{
		Link:           stat.task.Link,
		Filename:       pointer.To(stat.filename),
		WorkStatus:     entity.WORK,
		Message:        pointer.To(stat.msg),
		TargetQuantity: stat.task.Quality,
	}); err != nil {
		fmt.Println("statUpdate:", err.Error())
	}

	if err := d.cache.SetStatus(&temporary.TaskRequest{
		FileName:     stat.filename,
		Link:         stat.task.Link,
		MoveTo:       d.WorkDir,
		MaxQuality:   stat.task.Quality,
		Procentage:   stat.procentage,
		Status:       stat.status,
		DownloadSize: stat.totalSize,
		CurrentSize:  stat.currectSize,
		Message:      stat.msg,
	}); err != nil {
		fmt.Println("SetStatus:", err.Error())
	}
}

func (d *DownloaderSource) Downloader(task *Task) error {
	if _, err := d.sqlRepo.UpdateStatus(task.Link, entity.WORK); err != nil {
		fmt.Println("Downloader.UpdateStatus:", err.Error())
	}

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
		Output("%(title)s [%(duration_string)s].%(ext)s").
		NoRestrictFilenames().
		Fixup(ytdlp.FixupForce).
		Retries("10").
		NoWriteSubs().
		IgnoreErrors().
		IgnoreNoFormatsError().
		NoAbortOnError().
		RmCacheDir().
		EmbedThumbnail().
		DownloaderArgs("-S language \"ru,en\"").
		ProgressFunc(time.Duration(time.Millisecond*750), func(update ytdlp.ProgressUpdate) {
			size = (float64(update.DownloadedBytes) / 1024) / 1024 // К мегабайтам
			totalSize = (float64(update.TotalBytes) / 1024) / 1024 // К мегабайтам
			fmt.Println(update.PercentString(), fmt.Sprintf("[%.2f/%.2f]mb", size, totalSize), fmt.Sprintf("[%s]", update.Info.Format), update.Filename)

			status := string(update.Status)
			if strings.Contains(status, "finished") {
				status = "converting"
			}

			if filename != *update.Info.Filename {
				filename = *update.Info.Filename
				if err := d.sqlRepo.Update(&persistent.LinkModelRequest{
					Link:           task.Link,
					Filename:       pointer.To(filename),
					WorkStatus:     entity.WORK,
					Message:        pointer.To(fmt.Sprintf("%s [%s]", status, update.Info.Format)),
					TargetQuantity: int(pointer.Get(update.Info.Height)),
				}); err != nil {
					fmt.Println("Update:", err.Error())
				}
				task.Quality = int(pointer.Get(update.Info.Height))
			}

			d.statUpdate(statInfo{
				task:        task,
				msg:         fmt.Sprintf("%s [%s]", status, update.Info.Format),
				filename:    filename,
				totalSize:   totalSize,
				currectSize: size,
				procentage:  update.Percent(),
				status:      entity.WORK,
			})
			if d.eagleMode && update.Percent() > float64(d.PercentToNext) {
				toNext.Do(func() {
					<-d.workersPool
				})
			}

		})

	var wg sync.WaitGroup

	defer func() {
		wg.Wait()
		if !d.eagleMode {
			<-d.workersPool
		}
		fmt.Printf("%s [%d] finished\n\n", task.Link, task.Quality)
	}()

	if err := func() error {
		defer func() {
			d.cache.LinkDone(task.Link)
		}()

		var err_resp error
		for i := range d.totalStages {
			stg := d.Stages[(i + 1)]

			dl.UnsetFormat()
			if stg.IsFormat {
				dl.Format(fmt.Sprintf("bv*[height<=%d]+ba", qualtiy))
			}

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
				if task != nil {
					_, err := dl.Run(context.TODO(), task.Link)
					if err != nil {
						d.statUpdate(statInfo{
							task:        task,
							msg:         fmt.Sprintf("download failed on stage #%d with retries on stage %d. Reason: %s", i+1, retry, err.Error()),
							filename:    filename,
							totalSize:   totalSize,
							currectSize: size,
							procentage:  0,
							status:      entity.WORK,
						})
						err_resp = err
						fmt.Printf("download failed: %s\n", err.Error())
						continue
					}

					// Скачивание и конвертация прошли успешно
					d.statUpdate(statInfo{
						task:        task,
						msg:         fmt.Sprintf("download complete on stage #%d witn retries on stage %d", i+1, retry),
						filename:    filename,
						totalSize:   totalSize,
						currectSize: size,
						procentage:  100,
						status:      entity.WORK,
					})

					return nil
				}
			}
		}

		return err_resp
	}(); err != nil {
		fmt.Println("Downloader:", err.Error())
		return fmt.Errorf("downloader: %w", err)
	}

	return nil
}

func (d *DownloaderSource) autoScale(ctx context.Context) {
	t := time.NewTicker(time.Hour * TIMEOUT_WORKERS)

	for {
		select {
		case <-t.C:
			if len(d.workersPool) == cap(d.workersPool) {
				<-d.workersPool
				fmt.Println("Queue is scaled")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (d *DownloaderSource) Processor(ctx context.Context) {
	fmt.Println("Workers pool:", cap(d.workersPool))

	go func() {
		fmt.Println("Launch auto scale")
		defer fmt.Println("Turn off auto scale")

		d.autoScale(ctx)
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("===UPDATE YT-DLP===")

		if _, err := ytdlp.New().Update(context.Background()); err != nil {
			fmt.Println("Processor.UpdateStatus:", err.Error())
		}
	}()
	wg.Wait()

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
				if task.Link == "" {
					<-d.workersPool
					return
				}

				fmt.Println("New download:", task.Link)
				if err := d.Downloader(task); err != nil {
					fmt.Printf("downloader: %s\n", err.Error())
					// Помещаем обратно в пул
					<-d.workersPool

					if _, err := d.sqlRepo.UpdateStatus(task.Link, entity.NEW); err != nil {
						fmt.Println("Processor.UpdateStatus:", err.Error())
					}
					d.cache.LinkDone(task.Link)

					return
				}

				if _, err := d.sqlRepo.UpdateStatus(task.Link, entity.TO_SEND); err != nil {
					fmt.Println("Processor.UpdateStatus:", err.Error())
				}

				time.Sleep(time.Second * TIMEOUT_WORKERS)
			}()
			time.Sleep(time.Second * time.Duration(rand.IntN(TIMEOUT_WORKERS)+1))
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

	if _, err := d.sqlRepo.UpdateStatus(link.Link, entity.WORK); err != nil {
		fmt.Println("UpdateStatus:", err.Error())
	}

	return &Task{
		Link:    link.Link,
		Quality: link.TargetQuantity,
	}, nil
}
