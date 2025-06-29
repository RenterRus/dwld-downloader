package v1

import (
	"context"
	proto "dwld-downloader/docs/proto/v1"
	"dwld-downloader/internal/controller/grpc/V1/response"
	"dwld-downloader/internal/usecase"
	"fmt"

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

func (v *V1) WorkQueue(ctx context.Context, in *emptypb.Empty) (*proto.WorkQueueResponse, error) {
	tasks, err := v.u.WorkQueue()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.WorkQueueResponse{
		LinksInWork: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) History(ctx context.Context, in *emptypb.Empty) (*proto.HistoryResponse, error) {
	tasks, err := v.u.History()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.HistoryResponse{
		History: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}
