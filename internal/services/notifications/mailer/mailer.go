package mailer

import (
	"context"
	"log"
)

type Message struct {
	To      string
	Subject string
	Body    string
}

type Mailer interface {
	Send(ctx context.Context, message Message) error
}

type LogMailer struct{}

func (LogMailer) Send(_ context.Context, message Message) error {
	log.Printf("email sent to=%s subject=%q", message.To, message.Subject)
	return nil
}
