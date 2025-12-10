package ftp

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/dwld-downloader/internal/entity"
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	"github.com/RenterRus/dwld-downloader/internal/usecase"
	v1 "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	TIMEOUT_LOAD_SEC = 7
)

type FTPSender struct {
	sender  v1.SenderClient
	sqlRepo persistent.SQLRepo
	cache   temporary.CacheRepo
}

func (f *FTPSender) Sender(ctx context.Context) {
	t := time.NewTicker(time.Second * TIMEOUT_LOAD_SEC)
	for {
		select {
		case <-t.C:
			link, err := f.sqlRepo.SelectOne(entity.TO_SEND)
			if err != nil {
				fmt.Printf("Sedner(ftp): %s\n", err)
			}
			if link == nil {
				break
			}

			if _, err := f.sender.ToQueue(ctx, &v1.ToQueueRequest{
				Link:          link.Link,
				Filename:      pointer.Get(link.Filename),
				UserName:      link.UserName,
				TargetQuality: int32(link.TargetQuantity),
			}); err != nil {
				fmt.Printf("Sedner(ftp(sendToQueue)): %s\n", err)
				break
			}

			if _, err := f.sqlRepo.Delete(link.Link); err != nil {
				fmt.Println("Sender.Delete:", err.Error())
			}
		case <-ctx.Done():
			fmt.Println("Sender(ftp): context done")
			return
		}
	}
}

func (f *FTPSender) CleanHistory(ctx context.Context) {
	if _, err := f.sender.CleanDone(ctx, &emptypb.Empty{}); err != nil {
		fmt.Println("CleanHistory:", err.Error())
	}
}

func (f *FTPSender) Status(ctx context.Context) ([]*usecase.OnWork, error) {
	status, err := f.sender.LoadStatus(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Status: %w", err)
	}

	resp := make([]*usecase.OnWork, len(status.LinksInWork))
	for _, v := range status.LinksInWork {
		resp = append(resp, &usecase.OnWork{
			Link:           v.Link,
			Filename:       v.Filename,
			MoveTo:         v.MoveTo,
			TargetQuantity: v.TargetQuantity,
			Procentage:     v.Procentage,
			Status:         v.Status,
			TotalSize:      v.TotalSize,
			CurrentSize:    v.CurrentSize,
			Message:        v.Message,
		})
	}

	return resp, nil
}

func (f *FTPSender) Queue(ctx context.Context) ([]*usecase.Task, error) {
	queue, err := f.sender.LoadQueue(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("Queue(Queue): %w", err)
	}

	resp := make([]*usecase.Task, 0, len(queue.Queue))
	for _, v := range queue.Queue {
		resp = append(resp, &usecase.Task{
			Link:        v.Link,
			MaxQuantity: v.TargetQuality,
			Status:      v.Status,
			Name:        v.Name,
			Message:     v.Message,
		})
	}

	return resp, nil
}
