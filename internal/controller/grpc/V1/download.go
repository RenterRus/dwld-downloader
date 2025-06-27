package v1

import (
	"context"
	v1 "dwld-downloader/docs/proto/v1"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) SetToQueue(ctx context.Context, in *v1.SetToQueueRequest) (*v1.SetToQueueResponse, error) {
	fmt.Println("SetToQueue")
	return nil, nil
}

func (v *V1) DeleteFromQueue(ctx context.Context, in *v1.DeleteFromQueueRequest) (*v1.DeleteFromQueueResponse, error) {
	fmt.Println("DeleteFromQueue")

	return nil, nil
}

func (v *V1) CleanHistory(ctx context.Context, in *emptypb.Empty) (*v1.CleanHistoryResponse, error) {
	fmt.Println("CleanHistory")

	return nil, nil
}

func (v *V1) WorkQueue(ctx context.Context, in *emptypb.Empty) (*v1.WorkQueueResponse, error) {
	fmt.Println("WorkQueue")

	return nil, nil
}

func (v *V1) History(ctx context.Context, in *emptypb.Empty) (*v1.HistoryResponse, error) {
	fmt.Println("History")

	return nil, nil
}
