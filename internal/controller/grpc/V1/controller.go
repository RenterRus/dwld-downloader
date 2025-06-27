package v1

import (
	v1 "dwld-downloader/docs/proto/v1"
	"dwld-downloader/internal/usecase"
	"log"
)

type V1 struct {
	v1.DownloaderServer

	u usecase.Downloader
	l log.Logger
}
