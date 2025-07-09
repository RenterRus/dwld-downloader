package persistent

import "github.com/RenterRus/dwld-downloader/internal/entity"

type LinkModel struct {
	Link           string  `sql:"link"`
	Filename       *string `sql:"filename"`
	WorkStatus     string  `sql:"work_status"`
	Message        *string `sql:"message"`
	TargetQuantity int     `sql:"target_quantity"`
}

type LinkModelRequest struct {
	Link           string        `sql:"link"`
	Filename       *string       `sql:"filename"`
	WorkStatus     entity.Status `sql:"work_status"`
	Message        *string       `sql:"message"`
	TargetQuantity int           `sql:"target_quantity"`
}

type SQLRepo interface {
	SelectHistory(withoutStatus *entity.Status) ([]LinkModel, error)
	Insert(link string, maxQuality int) ([]LinkModel, error)
	UpdateStatus(link string, status entity.Status) ([]LinkModel, error)
	Delete(link string) ([]LinkModel, error)
	DeleteHistory() ([]LinkModel, error)

	SelectOne(status entity.Status) (*LinkModel, error)
	Update(*LinkModelRequest) error
}
