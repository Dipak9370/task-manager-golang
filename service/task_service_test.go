package service

import (
	"context"
	"task-manager/models"
	"testing"
)

type MockRepo struct {
	CreatedTask *models.Task
}

func (m *MockRepo) Create(ctx context.Context, task *models.Task) error {
	m.CreatedTask = task
	return nil
}

func (m *MockRepo) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return &models.Task{
		ID:     id,
		UserID: "user1",
		Status: "pending",
	}, nil
}

func (m *MockRepo) Update(ctx context.Context, task *models.Task) error { return nil }
func (m *MockRepo) Delete(ctx context.Context, id string) error         { return nil }
func (m *MockRepo) GetAll(ctx context.Context, userID, role, status string, limit, offset int) ([]models.Task, error) {
	return []models.Task{}, nil
}

func TestCreateTask(t *testing.T) {
	mock := &MockRepo{}
	service := TaskService{Repo: mock}

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
	service := TaskService{Repo: mock}

	_, err := service.GetTask(context.Background(), "task1", "user2", "user")
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}

func TestGetTask_AdminAccess(t *testing.T) {
	mock := &MockRepo{}
	service := TaskService{Repo: mock}

	task, err := service.GetTask(context.Background(), "task1", "anyuser", "admin")
	if err != nil {
		t.Fatal(err)
	}

	if task.ID != "task1" {
		t.Errorf("wrong task returned")
	}
}
