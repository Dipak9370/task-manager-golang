package service

import (
	"context"
	"task-manager/models"
	"task-manager/repository"
	"task-manager/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo repository.IUserRepository
}

func (s *AuthService) Register(ctx context.Context, email, password, role string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hash),
		Role:     role,
	}
	return s.Repo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.Repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	return utils.GenerateToken(user.ID, user.Role)
}
