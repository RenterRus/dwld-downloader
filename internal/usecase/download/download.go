package download

import "dwld-downloader/internal/usecase"

type downlaoder struct {
}

func NewDownload() usecase.Downloader {
	return &downlaoder{}
}
