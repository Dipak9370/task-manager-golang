package repository

import (
	"context"
	"task-manager/models"
)

type ITaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, userID, role, status string, limit, offset int) ([]models.Task, error)
}
