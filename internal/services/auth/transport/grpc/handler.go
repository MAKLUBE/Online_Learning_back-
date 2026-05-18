package grpc

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"online-learning-platform/internal/shared/authz"
	"online-learning-platform/internal/shared/grpcerrors"
	"online-learning-platform/internal/shared/validate"
	"sync"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/auth/model"
	"online-learning-platform/internal/services/auth/repository"
	"online-learning-platform/internal/services/auth/repository/memory"
	"online-learning-platform/internal/services/auth/usecase"
	authpb "online-learning-platform/pkg/gen/auth"
)

type Handler struct {
	authpb.UnimplementedAuthServiceServer
	useCase *usecase.AuthUseCase
	mu      sync.RWMutex
	users   map[string]*authpb.UserAuthResponse // legacy fallback for tests only
	store   repository.UserRepository
	events  events.Publisher
}

func NewHandler(useCase *usecase.AuthUseCase) *Handler {
	mem := memory.NewAuthRepository()
	return &Handler{useCase: useCase, users: make(map[string]*authpb.UserAuthResponse), store: mem, events: events.NoopPublisher{}}
}

func NewHandlerWithUserStore(useCase *usecase.AuthUseCase, store repository.UserRepository) *Handler {
	if store == nil {
		store = memory.NewAuthRepository()
	}
	return &Handler{useCase: useCase, users: make(map[string]*authpb.UserAuthResponse), store: store, events: events.NoopPublisher{}}
}

func NewHandlerWithDependencies(useCase *usecase.AuthUseCase, store repository.UserRepository, publisher events.Publisher) *Handler {
	handler := NewHandlerWithUserStore(useCase, store)
	if publisher != nil {
		handler.events = publisher
	}
	return handler
}

func (h *Handler) RegisterUser(ctx context.Context, req *authpb.RegisterUserRequest) (*authpb.UserAuthResponse, error) {
	if !validate.Email(req.GetEmail()) {
		return nil, grpcerrors.InvalidArgument("valid email is required")
	}
	if !validate.MinLen(req.GetPassword(), 6) {
		return nil, grpcerrors.InvalidArgument("password must contain at least 6 characters")
	}
	if !validate.Required(req.GetDisplayName()) {
		return nil, grpcerrors.InvalidArgument("display_name is required")
	}
	if _, err := h.store.GetUserByEmail(ctx, req.GetEmail()); err == nil {
		return nil, grpcerrors.AlreadyExists("user with this email already exists")
	} else if !errors.Is(err, sql.ErrNoRows) && err.Error() != "user not found" {
		return nil, grpcerrors.Internal(err)
	}
	now := time.Now().UTC()
	user, err := h.store.CreateUser(ctx, model.User{
		ID:           "user-" + events.NewID(),
		Email:        req.GetEmail(),
		DisplayName:  req.GetDisplayName(),
		Role:         "student",
		PasswordHash: hashPassword(req.GetPassword()),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, grpcerrors.AlreadyExists("user with this email already exists")
		}
		return nil, grpcerrors.Internal(err)
	}
	_ = h.events.Publish(ctx, "user.registered", events.New("user.registered", map[string]any{
		"user_id": user.ID,
		"email":   user.Email,
		"name":    user.DisplayName,
	}))
	return toUserResponse(user), nil
}

func (h *Handler) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.TokenResponse, error) {
	if !validate.Email(req.GetEmail()) {
		return nil, grpcerrors.InvalidArgument("valid email is required")
	}
	if !validate.Required(req.GetPassword()) {
		return nil, grpcerrors.InvalidArgument("password is required")
	}
	user, err := h.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, grpcerrors.Unauthenticated("invalid email or password")
	}
	if user.PasswordHash != "" && user.PasswordHash != hashPassword(req.GetPassword()) {
		return nil, grpcerrors.Unauthenticated("invalid email or password")
	}
	tokens, err := h.useCase.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return &authpb.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresAt: tokens.ExpiresAt.Format(time.RFC3339)}, nil
}

func (h *Handler) LogoutUser(ctx context.Context, req *authpb.LogoutUserRequest) (*authpb.Empty, error) {
	if !validate.Required(req.GetRefreshToken()) {
		return nil, grpcerrors.InvalidArgument("refresh_token is required")
	}
	return &authpb.Empty{}, grpcerrors.Internal(h.useCase.Logout(ctx, req.GetRefreshToken()))
}

func (h *Handler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.TokenResponse, error) {
	if !validate.Required(req.GetRefreshToken()) {
		return nil, grpcerrors.InvalidArgument("refresh_token is required")
	}
	tokens, err := h.useCase.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, grpcerrors.Unauthenticated(err.Error())
	}
	return &authpb.TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresAt: tokens.ExpiresAt.Format(time.RFC3339)}, nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	if !validate.Required(req.GetAccessToken()) {
		return nil, grpcerrors.InvalidArgument("access_token is required")
	}
	userID, err := h.useCase.ValidateToken(ctx, req.GetAccessToken())
	if err != nil {
		return nil, grpcerrors.Unauthenticated(err.Error())
	}
	return &authpb.ValidateTokenResponse{Valid: true, UserId: userID, Role: "student"}, nil
}

func (h *Handler) VerifyEmail(_ context.Context, req *authpb.VerifyEmailRequest) (*authpb.UserAuthResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	user, err := h.store.GetUserByID(context.Background(), req.GetUserId())
	if err != nil {
		return nil, grpcerrors.NotFound("user not found")
	}
	user.Verified = true
	user.UpdatedAt = time.Now().UTC()
	updated, err := h.store.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toUserResponse(updated), nil
}

func (h *Handler) RequestPasswordReset(ctx context.Context, req *authpb.RequestPasswordResetRequest) (*authpb.PasswordResetResponse, error) {
	if !validate.Email(req.GetEmail()) {
		return nil, grpcerrors.InvalidArgument("valid email is required")
	}
	token, err := h.useCase.RequestPasswordReset(ctx, req.GetEmail())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	_ = h.events.Publish(ctx, "user.password_reset_requested", events.New("user.password_reset_requested", map[string]any{
		"email":       req.GetEmail(),
		"reset_token": token,
	}))
	return &authpb.PasswordResetResponse{ResetToken: token}, nil
}

func (h *Handler) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.Empty, error) {
	if !validate.Required(req.GetResetToken()) {
		return nil, grpcerrors.InvalidArgument("reset_token is required")
	}
	if !validate.MinLen(req.GetNewPassword(), 6) {
		return nil, grpcerrors.InvalidArgument("new_password must contain at least 6 characters")
	}
	return &authpb.Empty{}, grpcerrors.Internal(h.useCase.ResetPassword(ctx, req.GetResetToken(), req.GetNewPassword()))
}

func (h *Handler) GetUserProfile(_ context.Context, req *authpb.GetUserProfileRequest) (*authpb.UserAuthResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	user, err := h.store.GetUserByID(context.Background(), req.GetUserId())
	if err != nil {
		return nil, grpcerrors.NotFound("user not found")
	}
	return toUserResponse(user), nil
}

func (h *Handler) UpdateUserProfile(_ context.Context, req *authpb.UpdateUserProfileRequest) (*authpb.UserAuthResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	if !validate.Required(req.GetDisplayName()) {
		return nil, grpcerrors.InvalidArgument("display_name is required")
	}
	user, err := h.store.GetUserByID(context.Background(), req.GetUserId())
	if err != nil {
		return nil, grpcerrors.NotFound("user not found")
	}
	user.DisplayName = req.GetDisplayName()
	user.UpdatedAt = time.Now().UTC()
	updated, err := h.store.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toUserResponse(updated), nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *authpb.DeleteUserRequest) (*authpb.Empty, error) {
	if err := authz.RequireRole(ctx, "admin"); err != nil {
		// DeleteUser is an admin action. Use metadata x-user-role: admin in Postman.
		return nil, err
	}
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	if _, err := h.store.GetUserByID(context.Background(), req.GetUserId()); err != nil {
		return nil, grpcerrors.NotFound("user not found")
	}
	return &authpb.Empty{}, grpcerrors.Internal(h.store.DeleteUser(context.Background(), req.GetUserId()))
}

func (h *Handler) ChangeUserRole(ctx context.Context, req *authpb.ChangeUserRoleRequest) (*authpb.UserAuthResponse, error) {
	if err := authz.RequireRole(ctx, "admin"); err != nil {
		// ChangeUserRole is an admin action. Use metadata x-user-role: admin in Postman.
		return nil, err
	}
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	if req.GetRole() != "student" && req.GetRole() != "instructor" && req.GetRole() != "admin" {
		return nil, grpcerrors.InvalidArgument("role must be student, instructor or admin")
	}
	user, err := h.store.GetUserByID(context.Background(), req.GetUserId())
	if err != nil {
		return nil, grpcerrors.NotFound("user not found")
	}
	user.Role = req.GetRole()
	user.UpdatedAt = time.Now().UTC()
	updated, err := h.store.UpdateUser(context.Background(), user)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toUserResponse(updated), nil
}

func toUserResponse(user model.User) *authpb.UserAuthResponse {
	return &authpb.UserAuthResponse{Id: user.ID, Email: user.Email, DisplayName: user.DisplayName, Role: user.Role, Verified: user.Verified}
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "unique constraint"))
}

func contains(value, part string) bool {
	return len(value) >= len(part) && (value == part || len(part) == 0 || stringContains(value, part))
}

func stringContains(value, part string) bool {
	for i := 0; i+len(part) <= len(value); i++ {
		if value[i:i+len(part)] == part {
			return true
		}
	}
	return false
}
