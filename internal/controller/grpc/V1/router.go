package v1

import (
	v1 "dwld-downloader/docs/proto/v1"

	"dwld-downloader/internal/usecase"
	"log"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases usecase.Downloader, l log.Logger) {
	r := &V1{
		u: usecases,
		l: l,
	}

	v1.RegisterDownloaderServer(app, r)
}
