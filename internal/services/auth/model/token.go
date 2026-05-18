package model

import "time"

type User struct {
	ID           string
	Email        string
	DisplayName  string
	Role         string
	PasswordHash string
	Verified     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type Session struct {
	UserID       string
	RefreshToken string
	ExpiresAt    time.Time
}
