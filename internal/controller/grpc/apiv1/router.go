package apiv1

import (
	proto "dwld-downloader/docs/proto/v1"

	"dwld-downloader/internal/usecase"

	pbgrpc "google.golang.org/grpc"
)

func NewDownloadRoutes(app *pbgrpc.Server, usecases usecase.Downloader) {
	r := &V1{
		u: usecases,
	}

	proto.RegisterDownloaderServer(app, r)
}
