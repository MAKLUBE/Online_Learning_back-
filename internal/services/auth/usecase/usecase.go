package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/auth/model"
	"online-learning-platform/internal/services/auth/repository"
)

var ErrCredentialsRequired = errors.New("email and password are required")

type AuthUseCase struct {
	repo repository.AuthRepository
}

func NewAuthUseCase(repo repository.AuthRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (u *AuthUseCase) Login(ctx context.Context, email, password string) (model.TokenPair, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return model.TokenPair{}, ErrCredentialsRequired
	}
	tokens := model.TokenPair{
		AccessToken:  tokenFor("access", email, 15*time.Minute),
		RefreshToken: tokenFor("refresh", email, 7*24*time.Hour),
		ExpiresAt:    time.Now().UTC().Add(15 * time.Minute),
	}
	if err := u.repo.SaveSession(ctx, model.Session{
		UserID:       strings.ToLower(strings.TrimSpace(email)),
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().UTC().Add(7 * 24 * time.Hour),
	}); err != nil {
		return model.TokenPair{}, err
	}
	return tokens, nil
}

func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	return u.repo.DeleteSession(ctx, refreshToken)
}

func (u *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (model.TokenPair, error) {
	session, err := u.repo.GetSession(ctx, refreshToken)
	if err != nil {
		return model.TokenPair{}, err
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		_ = u.repo.DeleteSession(ctx, refreshToken)
		return model.TokenPair{}, errors.New("refresh token expired")
	}
	return model.TokenPair{
		AccessToken:  tokenFor("access", session.UserID, 15*time.Minute),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().UTC().Add(15 * time.Minute),
	}, nil
}

func (u *AuthUseCase) ValidateToken(_ context.Context, accessToken string) (string, error) {
	if strings.TrimSpace(accessToken) == "" {
		return "", errors.New("token is required")
	}
	return "validated-user", nil
}

func (u *AuthUseCase) RequestPasswordReset(_ context.Context, email string) (string, error) {
	if strings.TrimSpace(email) == "" {
		return "", ErrCredentialsRequired
	}
	return tokenFor("reset", email, 30*time.Minute), nil
}

func (u *AuthUseCase) ResetPassword(_ context.Context, resetToken, newPassword string) error {
	if strings.TrimSpace(resetToken) == "" || strings.TrimSpace(newPassword) == "" {
		return ErrCredentialsRequired
	}
	return nil
}

func tokenFor(kind, subject string, ttl time.Duration) string {
	raw := kind + ":" + strings.ToLower(strings.TrimSpace(subject)) + ":" + time.Now().UTC().Add(ttl).Format(time.RFC3339Nano) + ":" + events.NewID()
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
