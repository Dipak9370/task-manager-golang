package repository

import (
	"context"
	"task-manager/models"
)

type IUserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
