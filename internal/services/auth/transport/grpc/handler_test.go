package grpc

import (
	"context"
	"testing"

	"online-learning-platform/internal/services/auth/repository/memory"
	"online-learning-platform/internal/services/auth/usecase"
	authpb "online-learning-platform/pkg/gen/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegisterUserValidation(t *testing.T) {
	handler := NewHandler(usecase.NewAuthUseCase(memory.NewAuthRepository()))

	_, err := handler.RegisterUser(context.Background(), &authpb.RegisterUserRequest{
		Email:       "wrong-email",
		Password:    "",
		DisplayName: "Student",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument, got %v", status.Code(err))
	}
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	handler := NewHandler(usecase.NewAuthUseCase(memory.NewAuthRepository()))
	req := &authpb.RegisterUserRequest{
		Email:       "student@example.com",
		Password:    "123456",
		DisplayName: "Student",
	}
	if _, err := handler.RegisterUser(context.Background(), req); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	_, err := handler.RegisterUser(context.Background(), req)
	if status.Code(err) != codes.AlreadyExists {
		t.Fatalf("expected AlreadyExists, got %v", status.Code(err))
	}
}
