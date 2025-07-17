package ftp

import (
	"context"

	"github.com/RenterRus/dwld-downloader/internal/usecase"
)

type File struct {
	Link          string
	Filename      string
	TargetQuality int
}
type Sender interface {
	Sender(ctx context.Context)
	CleanHistory(ctx context.Context)
	Status(ctx context.Context) ([]*usecase.OnWork, error)
	Queue(ctx context.Context) ([]*usecase.Task, error)
}
