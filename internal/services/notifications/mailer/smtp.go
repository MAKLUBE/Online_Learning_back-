package mailer

import (
	"context"
	"net/smtp"
	"os"
	"strings"
)

type SMTPMailer struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPFromEnv() SMTPMailer {
	return SMTPMailer{
		Host:     env("SMTP_HOST", "smtp.gmail.com"),
		Port:     env("SMTP_PORT", "587"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     env("SMTP_FROM", os.Getenv("SMTP_USERNAME")),
	}
}

func (m SMTPMailer) Send(ctx context.Context, message Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if m.Username == "" || m.Password == "" {
		return LogMailer{}.Send(ctx, message)
	}
	addr := m.Host + ":" + m.Port
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	body := strings.Join([]string{
		"From: " + m.From,
		"To: " + message.To,
		"Subject: " + message.Subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		message.Body,
	}, "\r\n")
	return smtp.SendMail(addr, auth, m.From, []string{message.To}, []byte(body))
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
