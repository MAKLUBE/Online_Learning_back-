package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	platformredis "online-learning-platform/internal/adapters/redis"
	"online-learning-platform/internal/services/auth/model"
	"online-learning-platform/internal/services/auth/repository"
)

type Repository struct {
	client *platformredis.Client
}

func New(client *platformredis.Client) repository.AuthRepository {
	return &Repository{client: client}
}

func (r *Repository) SaveSession(ctx context.Context, session model.Session) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		ttl = time.Hour
	}
	return r.client.Set(ctx, key(session.RefreshToken), string(payload), ttl)
}

func (r *Repository) DeleteSession(ctx context.Context, refreshToken string) error {
	return r.client.Del(ctx, key(refreshToken))
}

func (r *Repository) GetSession(ctx context.Context, refreshToken string) (model.Session, error) {
	payload, err := r.client.Get(ctx, key(refreshToken))
	if errors.Is(err, platformredis.ErrNil) {
		return model.Session{}, errors.New("session not found")
	}
	if err != nil {
		return model.Session{}, err
	}
	var session model.Session
	if err := json.Unmarshal([]byte(payload), &session); err != nil {
		return model.Session{}, err
	}
	return session, nil
}

func (r *Repository) Login(context.Context, string, string) (model.TokenPair, error) {
	return model.TokenPair{}, errors.New("login is handled in usecase")
}

func key(refreshToken string) string {
	return "auth:refresh:" + refreshToken
}
