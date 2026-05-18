package memory

import (
	"context"
	"errors"
	"sync"

	"online-learning-platform/internal/services/notifications/model"
)

type Repository struct {
	mu            sync.RWMutex
	notifications map[string]model.Notification
	settings      map[string]model.Settings
}

func New() *Repository {
	return &Repository{
		notifications: make(map[string]model.Notification),
		settings:      make(map[string]model.Settings),
	}
}

func (r *Repository) Create(_ context.Context, notification model.Notification) (model.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.notifications {
		if existing.IdempotencyKey == notification.IdempotencyKey && existing.Status == "sent" {
			return existing, nil
		}
	}
	r.notifications[notification.ID] = notification
	return notification, nil
}

func (r *Repository) GetByID(_ context.Context, id string) (model.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	notification, ok := r.notifications[id]
	if !ok {
		return model.Notification{}, errors.New("notification not found")
	}
	return notification, nil
}

func (r *Repository) ListByUser(_ context.Context, userID string) ([]model.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Notification, 0)
	for _, notification := range r.notifications {
		if notification.UserID == userID {
			items = append(items, notification)
		}
	}
	return items, nil
}

func (r *Repository) Update(_ context.Context, notification model.Notification) (model.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notifications[notification.ID] = notification
	return notification, nil
}

func (r *Repository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.notifications, id)
	return nil
}

func (r *Repository) UnreadCount(_ context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, notification := range r.notifications {
		if notification.UserID == userID && !notification.Read {
			count++
		}
	}
	return count, nil
}

func (r *Repository) GetSettings(_ context.Context, userID string) (model.Settings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	settings, ok := r.settings[userID]
	if !ok {
		return model.Settings{UserID: userID, EmailEnabled: true}, nil
	}
	return settings, nil
}

func (r *Repository) UpdateSettings(_ context.Context, settings model.Settings) (model.Settings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings[settings.UserID] = settings
	return settings, nil
}
