package apiv1

import (
	"context"
	proto "dwld-downloader/docs/proto/v1"
	"dwld-downloader/internal/controller/grpc/apiv1/response"
	"dwld-downloader/internal/usecase"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) SetToQueue(ctx context.Context, in *proto.SetToQueueRequest) (*proto.SetToQueueResponse, error) {
	tasks, err := v.u.SetToQueue(in.Link, *in.MaxQuality)
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.SetToQueueResponse{
		LinksInWork: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) DeleteFromQueue(ctx context.Context, in *proto.DeleteFromQueueRequest) (*proto.DeleteFromQueueResponse, error) {
	tasks, err := v.u.DeleteFromQueue(in.Link)
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.DeleteFromQueueResponse{
		LinksInWork: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) CleanHistory(ctx context.Context, in *emptypb.Empty) (*proto.CleanHistoryResponse, error) {
	tasks, err := v.u.CleanHistory()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.CleanHistoryResponse{
		History: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) Status(ctx context.Context, in *emptypb.Empty) (*proto.StatusResponse, error) {
	tasks, err := v.u.Status()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.StatusResponse{
		Sensors: tasks.Sensors,
		LinksInWork: lo.Map(tasks.LinksInWork, func(t *usecase.OnWork, _ int) *proto.OnWork {
			return &proto.OnWork{
				Link:           t.Link,
				Filename:       t.Filename,
				MoveTo:         t.MoveTo,
				TargetQuantity: t.TargetQuantity,
				Procentage:     t.Procentage,
				Status:         t.Status,
				TotalSize:      t.TotalSize,
				CurrentSize:    t.CurrentSize,
				Message:        t.Message,
			}
		}),
	}, nil
}

func (v *V1) Queue(ctx context.Context, in *emptypb.Empty) (*proto.HistoryResponse, error) {
	tasks, err := v.u.Queue()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.HistoryResponse{
		History: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) Healtheck(ctx context.Context, in *emptypb.Empty) (*proto.HealtheckResponse, error) {
	return &proto.HealtheckResponse{
		Message: pointer.To("OK"),
	}, nil
}
