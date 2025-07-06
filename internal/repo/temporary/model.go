package temporary

import (
	"context"
	"dwld-downloader/internal/entity"
)

type TaskRequest struct {
	FileName   string        `json:"filename"`
	Link       string        `json:"link"`
	MoveTo     string        `json:"move_to"`
	MaxQuality int           `json:"max_quantity"`
	Procentage float32       `json:"procentage"`
	Status     entity.Status `json:"status"`
	StageNum   int           `json:"stage_num"`
	StageConf  string        `json:"stage_conf"`
}

type TaskResp struct {
	MoveTo     string  `json:"move_to"`
	MaxQuality int     `json:"max_quantity"`
	Procentage float32 `json:"procentage"`
	Status     string  `json:"status"`
	StageNum   int     `json:"stage_num"`
	StageConf  string  `json:"stage_conf"`
}

type CacheResponse struct {
	//				link	 filename
	WorkStatus map[string]map[string]TaskResp
	Sensors    string
}

type CacheRepo interface {
	GetStatus() (*CacheResponse, error)
	SetStatus(*TaskRequest) error
	LinkDone(link string)
	Revisor(context.Context)
}
