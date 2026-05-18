package repository

import (
	"context"

	"online-learning-platform/internal/services/auth/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	GetUserByID(ctx context.Context, id string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, id string) error
}
