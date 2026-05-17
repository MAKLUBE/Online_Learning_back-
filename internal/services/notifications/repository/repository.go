package repository

import (
	"context"

	"online-learning-platform/internal/services/notifications/model"
)

type Repository interface {
	Create(ctx context.Context, notification model.Notification) (model.Notification, error)
	GetByID(ctx context.Context, id string) (model.Notification, error)
	ListByUser(ctx context.Context, userID string) ([]model.Notification, error)
	Update(ctx context.Context, notification model.Notification) (model.Notification, error)
	Delete(ctx context.Context, id string) error
	UnreadCount(ctx context.Context, userID string) (int, error)
	GetSettings(ctx context.Context, userID string) (model.Settings, error)
	UpdateSettings(ctx context.Context, settings model.Settings) (model.Settings, error)
}
