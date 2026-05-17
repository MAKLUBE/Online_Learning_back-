package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/notifications/mailer"
	"online-learning-platform/internal/services/notifications/model"
	"online-learning-platform/internal/services/notifications/repository"
)

type UseCase struct {
	repo   repository.Repository
	mailer mailer.Mailer
}

func New(repo repository.Repository, sender mailer.Mailer) *UseCase {
	if sender == nil {
		sender = mailer.LogMailer{}
	}
	return &UseCase{repo: repo, mailer: sender}
}

func (u *UseCase) SendEmail(ctx context.Context, userID, email, subject, body, notificationType string) (model.Notification, error) {
	if email == "" {
		return model.Notification{}, errors.New("recipient email is required")
	}
	notification := model.Notification{
		ID:             events.NewID(),
		UserID:         userID,
		RecipientEmail: email,
		Subject:        subject,
		Body:           body,
		Type:           notificationType,
		Status:         "pending",
		IdempotencyKey: idempotencyKey(userID, email, subject, body),
		CreatedAt:      time.Now().UTC(),
	}
	created, err := u.repo.Create(ctx, notification)
	if err != nil {
		return model.Notification{}, err
	}
	if err := u.mailer.Send(ctx, mailer.Message{To: email, Subject: subject, Body: body}); err != nil {
		created.Status = "failed"
		created.ErrorMessage = err.Error()
		updated, _ := u.repo.Update(ctx, created)
		return updated, err
	}
	created.Status = "sent"
	created.SentAt = time.Now().UTC()
	return u.repo.Update(ctx, created)
}

func (u *UseCase) SendWelcomeEmail(ctx context.Context, userID, email, name string) (model.Notification, error) {
	return u.SendEmail(ctx, userID, email, "Welcome to Online Learning", fmt.Sprintf("Hi %s, welcome to the adapters.", name), TypeWelcome)
}

func (u *UseCase) SendPasswordResetEmail(ctx context.Context, userID, email, resetToken string) (model.Notification, error) {
	return u.SendEmail(ctx, userID, email, "Password reset", "Use this reset token: "+resetToken, TypePasswordReset)
}

func (u *UseCase) SendCourseEnrollmentEmail(ctx context.Context, userID, email, courseID string) (model.Notification, error) {
	return u.SendEmail(ctx, userID, email, "Course enrollment", "You were enrolled in course "+courseID, TypeCourseEnrollment)
}

func (u *UseCase) SendAssignmentGradedEmail(ctx context.Context, userID, email, assignmentID string, grade float64) (model.Notification, error) {
	return u.SendEmail(ctx, userID, email, "Assignment graded", fmt.Sprintf("Assignment %s grade: %.2f", assignmentID, grade), TypeAssignmentGraded)
}

func (u *UseCase) SendDeadlineReminder(ctx context.Context, userID, email, assignmentID string) (model.Notification, error) {
	return u.SendEmail(ctx, userID, email, "Deadline reminder", "Assignment deadline is close: "+assignmentID, TypeDeadlineReminder)
}

func (u *UseCase) GetNotificationHistory(ctx context.Context, userID string) ([]model.Notification, error) {
	return u.repo.ListByUser(ctx, userID)
}

func (u *UseCase) GetNotificationByID(ctx context.Context, id string) (model.Notification, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) MarkNotificationRead(ctx context.Context, id string) (model.Notification, error) {
	notification, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return model.Notification{}, err
	}
	notification.Read = true
	return u.repo.Update(ctx, notification)
}

func (u *UseCase) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	return u.repo.UnreadCount(ctx, userID)
}

func (u *UseCase) UpdateNotificationSettings(ctx context.Context, userID string, enabled bool) (model.Settings, error) {
	return u.repo.UpdateSettings(ctx, model.Settings{UserID: userID, EmailEnabled: enabled})
}

func (u *UseCase) DeleteNotification(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *UseCase) ScheduleNotification(ctx context.Context, userID, email, subject, body string, _ time.Time) (model.Notification, error) {
	notification, err := u.SendEmail(ctx, userID, email, subject, body, "scheduled")
	if err != nil {
		return notification, err
	}
	notification.Status = "scheduled"
	return u.repo.Update(ctx, notification)
}

func idempotencyKey(parts ...string) string {
	raw := ""
	for _, part := range parts {
		raw += part + "|"
	}
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
