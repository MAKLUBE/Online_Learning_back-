package model

import "time"

type Notification struct {
	ID             string
	UserID         string
	RecipientEmail string
	Subject        string
	Body           string
	Type           string
	Status         string
	Read           bool
	IdempotencyKey string
	ErrorMessage   string
	CreatedAt      time.Time
	SentAt         time.Time
}

type Settings struct {
	UserID       string
	EmailEnabled bool
}

type Template struct {
	Name    string
	Subject string
	Body    string
}
