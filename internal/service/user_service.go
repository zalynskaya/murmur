package service

import (
	"context"

	custom_error "github.com/zalynskaya/murmur/internal/middleware"

	"github.com/zalynskaya/murmur/internal/entity"
	"github.com/zalynskaya/murmur/internal/repo"
)

type Service struct {
	storage *repo.UserStorage
}

type CreateUserDTO struct {
	Username string `json:"username"`
}

func NewUserService(storage *repo.UserStorage) *Service {
	return &Service{storage: storage}
}

func (u Service) Create(ctx context.Context, dto CreateUserDTO) (string, error) {
	exists, err := u.storage.IsExistsByUsername(ctx, dto.Username)
	if err != nil {
		return "", err
	}

	if exists {
		return "", custom_error.ErrUserDuplicate
	}

	user := entity.User{
		Username: dto.Username,
	}

	return u.storage.Create(ctx, user)
}
