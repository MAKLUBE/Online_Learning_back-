package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"online-learning-platform/internal/services/auth/model"
)

type AuthRepository struct {
	mu       sync.RWMutex
	sessions map[string]model.Session
	users    map[string]model.User
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{sessions: make(map[string]model.Session), users: make(map[string]model.User)}
}

func (r *AuthRepository) SaveSession(_ context.Context, session model.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[session.RefreshToken] = session
	return nil
}

func (r *AuthRepository) DeleteSession(_ context.Context, refreshToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, refreshToken)
	return nil
}

func (r *AuthRepository) GetSession(_ context.Context, refreshToken string) (model.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	session, ok := r.sessions[refreshToken]
	if !ok {
		return model.Session{}, errors.New("session not found")
	}
	return session, nil
}

func (r *AuthRepository) Login(ctx context.Context, email, password string) (model.TokenPair, error) {
	return model.TokenPair{
		AccessToken:  "legacy-access-token",
		RefreshToken: "legacy-refresh-token",
	}, nil
}

func (r *AuthRepository) CreateUser(_ context.Context, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.users {
		if existing.Email == user.Email {
			return model.User{}, errors.New("user already exists")
		}
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = user.CreatedAt
	}
	r.users[user.ID] = user
	return user, nil
}

func (r *AuthRepository) GetUserByID(_ context.Context, id string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *AuthRepository) GetUserByEmail(_ context.Context, email string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (r *AuthRepository) UpdateUser(_ context.Context, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return model.User{}, errors.New("user not found")
	}
	user.UpdatedAt = time.Now().UTC()
	r.users[user.ID] = user
	return user, nil
}

func (r *AuthRepository) DeleteUser(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(r.users, id)
	return nil
}
