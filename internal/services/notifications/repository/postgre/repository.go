package postgre

import (
	"context"
	"database/sql"

	"online-learning-platform/internal/services/notifications/model"
	"online-learning-platform/internal/services/notifications/repository"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) repository.Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, n model.Notification) (model.Notification, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO notifications (id,user_id,recipient_email,subject,body,type,status,read,idempotency_key,error_message,created_at,sent_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, n.ID, n.UserID, n.RecipientEmail, n.Subject, n.Body, n.Type, n.Status, n.Read, n.IdempotencyKey, n.ErrorMessage, n.CreatedAt, n.SentAt)
	return n, err
}
func (r *Repository) GetByID(ctx context.Context, id string) (model.Notification, error) {
	var n model.Notification
	err := r.db.QueryRowContext(ctx, `SELECT id,user_id,recipient_email,subject,body,type,status,read,idempotency_key,error_message,created_at,sent_at FROM notifications WHERE id=$1`, id).Scan(&n.ID, &n.UserID, &n.RecipientEmail, &n.Subject, &n.Body, &n.Type, &n.Status, &n.Read, &n.IdempotencyKey, &n.ErrorMessage, &n.CreatedAt, &n.SentAt)
	return n, err
}
func (r *Repository) ListByUser(ctx context.Context, userID string) ([]model.Notification, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,user_id,recipient_email,subject,body,type,status,read,idempotency_key,error_message,created_at,sent_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Notification, 0)
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.RecipientEmail, &n.Subject, &n.Body, &n.Type, &n.Status, &n.Read, &n.IdempotencyKey, &n.ErrorMessage, &n.CreatedAt, &n.SentAt); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, rows.Err()
}
func (r *Repository) Update(ctx context.Context, n model.Notification) (model.Notification, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET status=$1,read=$2,error_message=$3,sent_at=$4 WHERE id=$5`, n.Status, n.Read, n.ErrorMessage, n.SentAt, n.ID)
	return n, err
}
func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM notifications WHERE id=$1`, id)
	return err
}
func (r *Repository) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM notifications WHERE user_id=$1 AND read=false`, userID).Scan(&count)
	return count, err
}
func (r *Repository) GetSettings(ctx context.Context, userID string) (model.Settings, error) {
	var s model.Settings
	err := r.db.QueryRowContext(ctx, `SELECT user_id,email_enabled FROM notification_settings WHERE user_id=$1`, userID).Scan(&s.UserID, &s.EmailEnabled)
	if err == sql.ErrNoRows {
		return model.Settings{UserID: userID, EmailEnabled: true}, nil
	}
	return s, err
}
func (r *Repository) UpdateSettings(ctx context.Context, s model.Settings) (model.Settings, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO notification_settings (user_id,email_enabled) VALUES ($1,$2) ON CONFLICT (user_id) DO UPDATE SET email_enabled=$2`, s.UserID, s.EmailEnabled)
	return s, err
}
