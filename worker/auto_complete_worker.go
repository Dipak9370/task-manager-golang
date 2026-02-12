package worker

import (
	"context"
	"log"
	"task-manager/repository"
	"time"
)

var TaskQueue = make(chan string, 100)

func StartWorker(ctx context.Context, repo *repository.TaskRepository, minutes int) {
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return
		case <-ticker.C:
			if err := repo.AutoCompleteTasks(minutes); err != nil {
				log.Println("worker error:", err)
			}
		}
	}
}
