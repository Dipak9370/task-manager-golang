package service

import (
	"context"
	"errors"
	"time"

	"task-manager/models"
	"task-manager/repository"
	"task-manager/worker"

	"github.com/google/uuid"
)

type TaskService struct {
	Repo repository.ITaskRepository
}

func NewTaskService(repo repository.ITaskRepository) *TaskService {
	return &TaskService{Repo: repo}
}

// Create Task
func (s *TaskService) CreateTask(ctx context.Context, title, desc, userID string) error {
	task := &models.Task{
		ID:          uuid.NewString(),
		Title:       title,
		Description: desc,
		Status:      "pending",
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.Repo.Create(ctx, task); err != nil {
		return err
	}

	// Push to background worker
	worker.TaskQueue <- task.ID
	return nil
}

// Get Task with authorization
func (s *TaskService) GetTask(ctx context.Context, id, userID, role string) (*models.Task, error) {
	task, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if role != "admin" && task.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return task, nil
}

// List Tasks
func (s *TaskService) ListTasks(
	ctx context.Context,
	userID, role, status string,
	page, limit int,
) ([]models.Task, error) {

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	return s.Repo.GetAll(ctx, userID, role, status, limit, offset)
}

// Delete Task with authorization
func (s *TaskService) DeleteTask(ctx context.Context, id, userID, role string) error {
	task, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role != "admin" && task.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.Repo.Delete(ctx, id)
}
