package services

import (
	"context"

	"github.com/AniketSrivastava1/recruit/backend/internal/models"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
)

type UserService struct {
	store *stores.UserStore
}

func NewUserService(store *stores.UserStore) *UserService {
	return &UserService{store: store}
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) error {
	return s.store.Create(ctx, u)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return s.store.GetByID(ctx, id)
}
