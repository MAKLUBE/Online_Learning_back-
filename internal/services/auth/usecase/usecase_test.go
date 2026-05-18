package usecase

import (
	"context"
	"testing"

	"online-learning-platform/internal/services/auth/repository/memory"
)

func TestLoginStoresRefreshSession(t *testing.T) {
	uc := NewAuthUseCase(memory.NewAuthRepository())

	tokens, err := uc.Login(context.Background(), "student@example.com", "secret")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("expected access and refresh tokens")
	}

	refreshed, err := uc.RefreshToken(context.Background(), tokens.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken != tokens.RefreshToken {
		t.Fatal("expected refreshed access token with same refresh token")
	}
}
