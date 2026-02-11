package service

import (
	"context"
	"task-manager/models"
	"testing"
)

type MockUserRepo struct {
	SavedUser *models.User
}

func (m *MockUserRepo) Create(ctx context.Context, user *models.User) error {
	m.SavedUser = user
	return nil
}

func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return m.SavedUser, nil
}

func TestRegister_HashesPassword(t *testing.T) {
	mock := &MockUserRepo{}
	service := AuthService{Repo: mock}

	err := service.Register(context.Background(), "test@gmail.com", "123456", "user")
	if err != nil {
		t.Fatal(err)
	}

	if mock.SavedUser.Password == "123456" {
		t.Fatal("password was not hashed")
	}
}

func TestLogin_GeneratesToken(t *testing.T) {
	mock := &MockUserRepo{}
	service := AuthService{Repo: mock}

	// First register user
	service.Register(context.Background(), "test@gmail.com", "123456", "user")

	token, err := service.Login(context.Background(), "test@gmail.com", "123456")
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("expected JWT token")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	mock := &MockUserRepo{}
	service := AuthService{Repo: mock}

	service.Register(context.Background(), "test@gmail.com", "123456", "user")

	_, err := service.Login(context.Background(), "test@gmail.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}
