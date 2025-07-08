package download

import (
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
	"dwld-downloader/internal/usecase"
	"fmt"

	"github.com/samber/lo"
)

type downlaoder struct {
	dbRepo    persistent.SQLRepo
	cacheRepo temporary.CacheRepo
}

func NewDownload(dbRepo persistent.SQLRepo, cache temporary.CacheRepo) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
	}
}

func (d *downlaoder) SetToQueue(link string, maxQuantity int32) ([]*usecase.Task, error) {
	resp, err := d.dbRepo.Insert(link, int(maxQuantity))
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return lo.Map(resp, LinkToTask), nil
}

func (d *downlaoder) DeleteFromQueue(link string) ([]*usecase.Task, error) {
	resp, err := d.dbRepo.Delete(link)
	if err != nil {
		return nil, fmt.Errorf("DeleteFromQueue: %w", err)
	}

	return lo.Map(resp, LinkToTask), nil
}

func (d *downlaoder) CleanHistory() ([]*usecase.Task, error) {
	resp, err := d.dbRepo.DeleteHistory()
	if err != nil {
		return nil, fmt.Errorf("CleanHistory: %w", err)
	}

	return lo.Map(resp, LinkToTask), nil
}

func (d *downlaoder) Status() (*usecase.StatusResponse, error) {
	resp, err := d.cacheRepo.GetStatus()
	if err != nil {
		return nil, fmt.Errorf("CleanHistory: %w", err)
	}

	links := make([]*usecase.OnWork, 0, len(resp.WorkStatus)*2)

	for link, v := range resp.WorkStatus {
		for file, info := range v {
			links = append(links, &usecase.OnWork{
				Link:           link,
				Filename:       file,
				MoveTo:         info.MoveTo,
				TargetQuantity: int64(info.MaxQuality),
				Procentage:     info.Procentage,
				Status:         info.Status,
				TotalSize:      info.DownloadSize,
				CurrentSize:    info.CurrentSize,
				Message:        info.Message,
			})
		}
	}

	return &usecase.StatusResponse{
		Sensors:     resp.Sensors,
		LinksInWork: links,
	}, nil
}

func (d *downlaoder) Queue() ([]*usecase.Task, error) {
	resp, err := d.dbRepo.SelectHistory(nil)
	if err != nil {
		return nil, fmt.Errorf("Queue: %w", err)
	}

	return lo.Map(resp, LinkToTask), nil
}
