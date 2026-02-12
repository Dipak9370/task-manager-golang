package repository

import (
	"context"
	"task-manager/models"
	"time"

	"gorm.io/gorm"
)

type TaskRepository struct {
	DB *gorm.DB
}

func (r *TaskRepository) AutoCompleteTasks(minutes int) error {
	cutoff := time.Now().Add(-time.Duration(minutes) * time.Minute)

	return r.DB.Model(&models.Task{}).
		Where("status = ? AND created_at <= ?", "pending", cutoff).
		Update("status", "completed").Error
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.DB.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	err := r.DB.WithContext(ctx).First(&task, "id = ?", id).Error
	return &task, err
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.DB.WithContext(ctx).Save(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	return r.DB.WithContext(ctx).Delete(&models.Task{}, "id = ?", id).Error
}

func (r *TaskRepository) GetAll(
	ctx context.Context,
	userID, role, status string,
	limit, offset int,
) ([]models.Task, error) {

	var tasks []models.Task
	query := r.DB.WithContext(ctx).Model(&models.Task{})

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Limit(limit).Offset(offset).Find(&tasks).Error
	return tasks, err
}
