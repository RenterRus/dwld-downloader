package usecase

type Stage struct {
	Configure string
	Progress  string
	Message   string
}

type Task struct {
	Link        string
	MaxQuantity string
	Status      string
	Name        *string
	Stage       *Stage
}

type Downloader interface {
	SetToQueue(link string, maxQuantity int32) ([]Task, error)
	DeleteFromQueue(link string) ([]Task, error)
	CleanHistory() ([]Task, error)
	WorkQueue() ([]Task, error)
	History() ([]Task, error)
}
