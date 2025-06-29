package download

import (
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/repo/temporary"
	"dwld-downloader/internal/usecase"
)

type downlaoder struct {
	dbRepo    *persistent.SQLRepo
	cacheRepo *temporary.Cache
}

func NewDownload(dbRepo *persistent.SQLRepo, cache *temporary.Cache) usecase.Downloader {
	return &downlaoder{
		dbRepo:    dbRepo,
		cacheRepo: cache,
	}
}
