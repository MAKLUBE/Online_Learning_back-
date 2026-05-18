package repository

import (
	"context"

	"online-learning-platform/internal/services/auth/model"
)

type AuthRepository interface {
	SaveSession(ctx context.Context, session model.Session) error
	DeleteSession(ctx context.Context, refreshToken string) error
	GetSession(ctx context.Context, refreshToken string) (model.Session, error)
	Login(ctx context.Context, email, password string) (model.TokenPair, error)
}
