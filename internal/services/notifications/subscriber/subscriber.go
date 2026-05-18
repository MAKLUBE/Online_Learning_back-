package subscriber

import (
	"context"
	"fmt"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/notifications/usecase"
)

type Handler struct {
	useCase *usecase.UseCase
}

func New(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Handle(ctx context.Context, event events.Event) error {
	userID := stringField(event, "user_id")
	email := stringField(event, "email")
	if email == "" {
		email = "student@example.com"
	}
	switch event.Type {
	case "user.registered":
		_, err := h.useCase.SendWelcomeEmail(ctx, userID, email, stringField(event, "name"))
		return err
	case "user.password_reset_requested":
		_, err := h.useCase.SendPasswordResetEmail(ctx, userID, email, stringField(event, "reset_token"))
		return err
	case "course.enrolled":
		_, err := h.useCase.SendCourseEnrollmentEmail(ctx, userID, email, stringField(event, "course_id"))
		return err
	case "assignment.graded":
		_, err := h.useCase.SendAssignmentGradedEmail(ctx, userID, email, stringField(event, "assignment_id"), floatField(event, "grade"))
		return err
	case "course.completed":
		_, err := h.useCase.SendEmail(ctx, userID, email, "Course completed", "You completed course "+stringField(event, "course_id"), "course_completed")
		return err
	default:
		return nil
	}
}

func stringField(event events.Event, key string) string {
	value, ok := event.Data[key]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprint(value)
}

func floatField(event events.Event, key string) float64 {
	value, ok := event.Data[key]
	if !ok || value == nil {
		return 0
	}
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}
