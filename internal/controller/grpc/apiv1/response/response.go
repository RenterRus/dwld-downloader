package response

import (
	proto "dwld-downloader/docs/proto/v1"
	"dwld-downloader/internal/usecase"
)

func TasksToLinks(task *usecase.Task) *proto.Task {
	return &proto.Task{
		Link:        task.Link,
		MaxQuantity: task.MaxQuantity,
		Status:      task.Status,
		Name:        task.Name,
		Message:     task.Message,
	}
}
