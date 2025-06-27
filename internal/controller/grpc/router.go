package grpc

import (
	v1 "dwld-downloader/internal/controller/grpc/v1"
	"dwld-downloader/internal/usecase"
	"log"

	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(app *pbgrpc.Server, usecases usecase.Downloader, l log.Logger) {
	v1.NewDownloadRoutes(app, usecases, l)
	reflection.Register(app)
}
