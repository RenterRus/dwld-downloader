package response

import (
	proto "github.com/RenterRus/dwld-downloader/docs/proto/v1"

	"github.com/RenterRus/dwld-downloader/internal/usecase"
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
