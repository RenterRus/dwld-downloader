package v1

import (
	"context"
	proto "dwld-downloader/docs/proto/v1"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (v *V1) SetToQueue(ctx context.Context, in *proto.SetToQueueRequest) (*proto.SetToQueueResponse, error) {
	fmt.Println("SetToQueue")
	return nil, nil
}

func (v *V1) DeleteFromQueue(ctx context.Context, in *proto.DeleteFromQueueRequest) (*proto.DeleteFromQueueResponse, error) {
	fmt.Println("DeleteFromQueue")

	return nil, nil
}

func (v *V1) CleanHistory(ctx context.Context, in *emptypb.Empty) (*proto.CleanHistoryResponse, error) {
	fmt.Println("CleanHistory")

	return nil, nil
}

func (v *V1) WorkQueue(ctx context.Context, in *emptypb.Empty) (*proto.WorkQueueResponse, error) {
	fmt.Println("WorkQueue")

	return nil, nil
}

func (v *V1) History(ctx context.Context, in *emptypb.Empty) (*proto.HistoryResponse, error) {
	fmt.Println("History")

	return nil, nil
}
