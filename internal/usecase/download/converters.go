package download

import (
	"dwld-downloader/internal/repo/persistent"
	"dwld-downloader/internal/usecase"
	"strconv"
)

func LinkToTask(item persistent.LinkModel, _ int) *usecase.Task {
	return &usecase.Task{
		Link:        item.Link,
		MaxQuantity: strconv.Itoa(item.TargetQuantity),
		Status:      item.WorkStatus,
		Name:        item.Filename,
		Message:     item.Message,
	}
}
