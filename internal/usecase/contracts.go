package usecase

type OnWork struct {
	Link           string
	Filename       string
	MoveTo         string
	TargetQuantity int64
	Procentage     float64
	Status         string
	TotalSize      float64
	CurrentSize    float64
	Message        string
}

type StatusResponse struct {
	Sensors     string
	LinksInWork []*OnWork
}
type Task struct {
	Link        string
	MaxQuantity string
	Status      string
	Name        *string
	Message     *string
}

type Downloader interface {
	SetToQueue(link string, maxQuantity int32) ([]*Task, error)
	DeleteFromQueue(link string) ([]*Task, error)
	CleanHistory() ([]*Task, error)
	Status() (*StatusResponse, error)
	Queue() ([]*Task, error)
}
