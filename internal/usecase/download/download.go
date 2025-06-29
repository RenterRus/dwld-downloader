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

func (d *downlaoder) SetToQueue(link string, maxQuantity int32) ([]*usecase.Task, error) {
	return nil, nil
}

func (d *downlaoder) DeleteFromQueue(link string) ([]*usecase.Task, error) {
	return nil, nil
}

func (d *downlaoder) CleanHistory() ([]*usecase.Task, error) {
	return nil, nil
}

func (d *downlaoder) WorkQueue() ([]*usecase.Task, error) {
	return nil, nil
}

func (d *downlaoder) History() ([]*usecase.Task, error) {
	return nil, nil
}
