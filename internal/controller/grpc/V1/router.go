package v1

import (
	v1 "dwld-downloader/docs/proto/v1"

	"dwld-downloader/internal/usecase"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases usecase.Downloader) {
	r := &V1{
		u: usecases,
	}

	v1.RegisterDownloaderServer(app, r)
}
