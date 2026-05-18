package usecase

import (
	"context"
	"testing"

	"online-learning-platform/internal/services/notifications/mailer"
	"online-learning-platform/internal/services/notifications/repository/memory"
)

type spyMailer struct {
	count int
}

func (m *spyMailer) Send(context.Context, mailer.Message) error {
	m.count++
	return nil
}

func TestSendWelcomeEmailSavesHistory(t *testing.T) {
	sender := &spyMailer{}
	uc := New(memory.New(), sender)

	notification, err := uc.SendWelcomeEmail(context.Background(), "user-1", "student@example.com", "Student")
	if err != nil {
		t.Fatalf("send failed: %v", err)
	}
	if notification.Status != "sent" {
		t.Fatalf("expected sent status, got %s", notification.Status)
	}
	if sender.count != 1 {
		t.Fatalf("expected one email, got %d", sender.count)
	}

	history, err := uc.GetNotificationHistory(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected one history item, got %d", len(history))
	}
}
