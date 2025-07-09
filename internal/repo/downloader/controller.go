package downloader

import "context"

type Downloader interface {
	Downloader(*Task) error
	Processor(context.Context)
	GetLink() (*Task, error)
}
