package download

import (
	"context"
	"fmt"

	"github.com/RenterRus/dwld-downloader/internal/controller/ftp"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	"github.com/RenterRus/dwld-downloader/internal/usecase"
	"github.com/samber/lo"
)

type downlaoder struct {
	dbRepo    persistent.SQLRepo
	cacheRepo temporary.CacheRepo
	ftpSender ftp.Sender
}

func NewDownload(dbRepo persistent.SQLRepo, cache temporary.CacheRepo, ftpSender ftp.Sender) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
		ftpSender: ftpSender,
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
	d.ftpSender.CleanHistory(context.Background())
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

	files, err := d.ftpSender.Status(context.Background())
	if err != nil {
		fmt.Printf("Status(ftp Status): %s\n", err.Error())
	} else {
		links = append(links, files...)
	}

	return &usecase.StatusResponse{
		Sensors:     resp.Sensors,
		LinksInWork: links,
	}, nil
}

func (d *downlaoder) Queue() ([]*usecase.Task, error) {
	queue, err := d.dbRepo.SelectHistory(nil)
	if err != nil {
		return nil, fmt.Errorf("Queue: %w", err)
	}

	resp := lo.Map(queue, LinkToTask)

	files, err := d.ftpSender.Queue(context.Background())
	if err != nil {
		fmt.Printf("Queue(ftp Queue): %s\n", err.Error())
	} else {
		resp = append(resp, files...)
	}

	return resp, nil
}
