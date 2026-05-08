package service

import (
	"context"
	"fluxera/internal/models"
	"fluxera/internal/repositories"
)

type UserService struct {
	users *repositories.UserRepository
}

func NewUserService(users *repositories.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.users.FindUserByID(ctx, id)
}
