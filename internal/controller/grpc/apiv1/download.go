package apiv1

import (
	"context"
	"fmt"

	"github.com/RenterRus/dwld-downloader/internal/controller/grpc/apiv1/response"

	proto "github.com/RenterRus/dwld-downloader/docs/proto/v1"

	"github.com/RenterRus/dwld-downloader/internal/usecase"

	"github.com/AlekSi/pointer"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) SetToQueue(ctx context.Context, in *proto.SetToQueueRequest) (*proto.SetToQueueResponse, error) {
	fmt.Println("================\nSetToQueue")
	defer func() {
		fmt.Println(in)
		fmt.Println("================")
	}()

	if in == nil || pointer.Get(in).Link == "" {
		return nil, fmt.Errorf("SetToQueue: empty request")
	}

	tasks, err := v.u.SetToQueue(in.GetLink(), *in.MaxQuality)
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
	fmt.Println("================\nDeleteFromQueue")
	defer func() {
		fmt.Println(in)
		fmt.Println("================")
	}()

	tasks, err := v.u.DeleteFromQueue(in.GetLink())
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
	fmt.Println("================\nCleanHistory")
	defer func() {
		fmt.Println("================")
	}()

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
	fmt.Println("================\nStatus")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Status()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.StatusResponse{
		Sensors: tasks.Sensors,
		LinksInWork: lo.FilterMap(tasks.LinksInWork, func(t *usecase.OnWork, _ int) (*proto.OnWork, bool) {
			if t == nil {
				return nil, false
			}
			return &proto.OnWork{
				Link:           t.Link,
				Filename:       t.Filename,
				MoveTo:         t.MoveTo,
				TargetQuantity: t.TargetQuantity,
				Procentage:     t.Procentage,
				Status:         t.Status,
				CurrentSize:    t.CurrentSize,
				TotalSize:      t.TotalSize,
				Message:        t.Message,
			}, true
		}),
	}, nil
}

func (v *V1) Queue(ctx context.Context, in *emptypb.Empty) (*proto.HistoryResponse, error) {
	fmt.Println("================\nQueue")
	defer func() {
		fmt.Println("================")
	}()

	tasks, err := v.u.Queue()
	if err != nil {
		return nil, fmt.Errorf("SetToQueue: %w", err)
	}

	return &proto.HistoryResponse{
		Queue: lo.Map(tasks, func(t *usecase.Task, _ int) *proto.Task {
			return response.TasksToLinks(t)
		}),
	}, nil
}

func (v *V1) Healtheck(ctx context.Context, in *emptypb.Empty) (*proto.HealtheckResponse, error) {
	fmt.Println("================\nHealtheck")
	defer func() {
		fmt.Println("================")
	}()

	return &proto.HealtheckResponse{
		Message: pointer.To("OK"),
	}, nil
}
