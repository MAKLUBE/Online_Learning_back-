package redis

import (
	"context"
	"errors"
	"time"
)

type IdempotencyStore struct {
	client *Client
}

func NewIdempotencyStore(client *Client) *IdempotencyStore {
	return &IdempotencyStore{client: client}
}

func (s *IdempotencyStore) AlreadyProcessed(ctx context.Context, key string) (bool, error) {
	value, err := s.client.Get(ctx, "idempotency:"+key)
	if errors.Is(err, ErrNil) {
		return false, nil
	}
	return value == "done", err
}

func (s *IdempotencyStore) MarkProcessed(ctx context.Context, key string, ttl time.Duration) error {
	return s.client.Set(ctx, "idempotency:"+key, "done", ttl)
}
