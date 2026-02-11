package worker

import (
	"context"
	"task-manager/repository"
	"time"
)

var TaskQueue = make(chan string, 100)

func StartWorker(repo *repository.TaskRepository, minutes int) {
	go func() {
		for taskID := range TaskQueue {
			time.Sleep(time.Duration(minutes) * time.Minute)

			ctx := context.Background()
			task, err := repo.GetByID(ctx, taskID)
			if err == nil && (task.Status == "pending" || task.Status == "in_progress") {
				task.Status = "completed"
				repo.Update(ctx, task)
			}
		}
	}()
}
