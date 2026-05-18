package postgre

import (
	"context"
	"database/sql"
	"errors"

	"online-learning-platform/internal/services/auth/model"
	"online-learning-platform/internal/services/auth/repository"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) repository.AuthRepository {
	return &Repository{db: db}
}

func (r *Repository) SaveSession(ctx context.Context, session model.Session) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO auth_sessions (refresh_token,user_id,expires_at) VALUES ($1,$2,$3)`, session.RefreshToken, session.UserID, session.ExpiresAt)
	return err
}

func (r *Repository) DeleteSession(ctx context.Context, refreshToken string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM auth_sessions WHERE refresh_token=$1`, refreshToken)
	return err
}

func (r *Repository) GetSession(ctx context.Context, refreshToken string) (model.Session, error) {
	var session model.Session
	err := r.db.QueryRowContext(ctx, `SELECT user_id,refresh_token,expires_at FROM auth_sessions WHERE refresh_token=$1`, refreshToken).
		Scan(&session.UserID, &session.RefreshToken, &session.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Session{}, errors.New("session not found")
	}
	return session, err
}

func (r *Repository) Login(context.Context, string, string) (model.TokenPair, error) {
	return model.TokenPair{}, errors.New("login is handled in usecase")
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id,email,display_name,role,password_hash,verified,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		user.ID, user.Email, user.DisplayName, user.Role, user.PasswordHash, user.Verified, user.CreatedAt, user.UpdatedAt)
	return user, err
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	return r.scanUser(ctx, `SELECT id,email,display_name,role,password_hash,verified,created_at,updated_at FROM users WHERE id=$1`, id)
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	return r.scanUser(ctx, `SELECT id,email,display_name,role,password_hash,verified,created_at,updated_at FROM users WHERE email=$1`, email)
}

func (r *Repository) UpdateUser(ctx context.Context, user model.User) (model.User, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET display_name=$1,role=$2,verified=$3,updated_at=$4 WHERE id=$5`,
		user.DisplayName, user.Role, user.Verified, user.UpdatedAt, user.ID)
	return user, err
}

func (r *Repository) DeleteUser(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id=$1`, id)
	return err
}

func (r *Repository) scanUser(ctx context.Context, query string, args ...any) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.DisplayName, &user.Role, &user.PasswordHash, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, errors.New("user not found")
	}
	return user, err
}
