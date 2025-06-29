package download

import (
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/usecase"
)

type downlaoder struct {
	dbRepo persistent.SQLRepo
}

func NewDownload(dbRepo persistent.SQLRepo) usecase.Downloader {
	return &downlaoder{
		dbRepo: dbRepo,
	}
}
