package persistent

import "dwld-downloader/internal/entity"

type LinkModel struct {
	Link           string  `sql:"link"`
	Filename       *string `sql:"filename"`
	Path           string  `sql:"path"`
	WorkStatus     string  `sql:"work_status"`
	StageConfig    *string `sql:"stage_config"`
	Message        *string `sql:"message"`
	Retry          *int    `sql:"retry"`
	TargetQuantity int     `sql:"target_quantity"`
}

type SQLRepo interface {
	SelectHistory() ([]LinkModel, error)
	Insert(link string, maxQuality int) ([]LinkModel, error)
	UpdateStatus(link string, status entity.Status) ([]LinkModel, error)
	Delete(link string) ([]LinkModel, error)
	DeleteHistory() ([]LinkModel, error)
}
