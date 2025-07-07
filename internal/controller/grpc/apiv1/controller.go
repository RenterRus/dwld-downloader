package apiv1

import (
	proto "dwld-downloader/docs/proto/v1"
	"dwld-downloader/internal/usecase"
)

type V1 struct {
	proto.DownloaderServer

	u usecase.Downloader
}
