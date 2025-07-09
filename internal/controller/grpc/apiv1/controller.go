package apiv1

import (
	proto "github.com/RenterRus/dwld-downloader/docs/proto/v1"
	"github.com/RenterRus/dwld-downloader/internal/usecase"
)

type V1 struct {
	proto.DownloaderServer

	u usecase.Downloader
}
