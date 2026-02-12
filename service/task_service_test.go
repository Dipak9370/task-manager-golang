package service

import (
	"context"
	"errors"
	"task-manager/models"
	"testing"
)

type MockRepo struct {
	CreatedTask *models.Task
}

// Delete implements [repository.ITaskRepository].
func (m *MockRepo) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// GetByID implements [repository.ITaskRepository].
func (m *MockRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	panic("unimplemented")
}

// Update implements [repository.ITaskRepository].
func (m *MockRepo) Update(ctx context.Context, task *models.Task) error {
	panic("unimplemented")
}

// Create
func (m *MockRepo) Create(ctx context.Context, task *models.Task) error {
	m.CreatedTask = task
	return nil
}

// Get with access control
func (m *MockRepo) GetByIDWithAccess(ctx context.Context, id, userID, role string) (*models.Task, error) {
	task := &models.Task{
		ID:     id,
		UserID: "user1",
		Status: "pending",
	}

	if role != "admin" && task.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return task, nil
}

// Delete with access control
func (m *MockRepo) DeleteWithAccess(ctx context.Context, id, userID, role string) (bool, error) {
	if role != "admin" && userID != "user1" {
		return false, nil
	}
	return true, nil
}

// List
func (m *MockRepo) GetAll(ctx context.Context, userID, role, status string, limit, offset int) ([]models.Task, error) {
	return []models.Task{}, nil
}

func TestCreateTask(t *testing.T) {
	mock := &MockRepo{}
	service := NewTaskService(mock)

	err := service.CreateTask(context.Background(), "Title", "Desc", "user1")
	if err != nil {
		t.Fatal(err)
	}

	if mock.CreatedTask == nil {
		t.Fatal("task was not created")
	}

	if mock.CreatedTask.Title != "Title" {
		t.Errorf("expected Title, got %s", mock.CreatedTask.Title)
	}

	if mock.CreatedTask.Status != "pending" {
		t.Errorf("expected pending status")
	}
}

func TestGetTask_Unauthorized(t *testing.T) {
	mock := &MockRepo{}
	service := NewTaskService(mock)

	_, err := service.GetTask(context.Background(), "task1", "user2", "user")
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}

func TestGetTask_AdminAccess(t *testing.T) {
	mock := &MockRepo{}
	service := NewTaskService(mock)

	task, err := service.GetTask(context.Background(), "task1", "anyuser", "admin")
	if err != nil {
		t.Fatal(err)
	}

	if task.ID != "task1" {
		t.Errorf("wrong task returned")
	}
}

func TestDeleteTask_Unauthorized(t *testing.T) {
	mock := &MockRepo{}
	service := NewTaskService(mock)

	err := service.DeleteTask(context.Background(), "task1", "user2", "user")
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}

func TestDeleteTask_Admin(t *testing.T) {
	mock := &MockRepo{}
	service := NewTaskService(mock)

	err := service.DeleteTask(context.Background(), "task1", "any", "admin")
	if err != nil {
		t.Fatal(err)
	}
}
