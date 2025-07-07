package temporary

import (
	"context"
	"dwld-downloader/internal/entity"
)

type TaskRequest struct {
	FileName     string        `json:"filename"`
	Link         string        `json:"link"`
	MoveTo       string        `json:"move_to"`
	MaxQuality   int           `json:"max_quantity"`
	Procentage   float64       `json:"procentage"`
	Status       entity.Status `json:"status"`
	DownloadSize float64       `json:"download_size"`
	CurrentSize  float64       `json:"current_size"`
	Message      string        `json:"message"`
}

type TaskResp struct {
	MoveTo       string  `json:"move_to"`
	MaxQuality   int     `json:"max_quantity"`
	Procentage   float64 `json:"procentage"`
	Status       string  `json:"status"`
	DownloadSize float64 `json:"download_size"`
	CurrentSize  float64 `json:"current_size"`
	Message      string  `json:"message"`
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
