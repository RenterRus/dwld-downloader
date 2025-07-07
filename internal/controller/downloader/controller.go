package downloader

type Downloader interface {
	Downloader(link string, quality int32)
	Processor()
	GetLink() (string, error)
}
