package usecase

import (
	"context"
	"testing"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/learning/repository/memory"
)

type capturedPublisher struct {
	subjects []string
}

func (p *capturedPublisher) Publish(_ context.Context, subject string, _ events.Event) error {
	p.subjects = append(p.subjects, subject)
	return nil
}

func TestMarkLessonCompletedPublishesCourseCompleted(t *testing.T) {
	publisher := &capturedPublisher{}
	uc := New(memory.New(), AlwaysExistsChecker{}, AlwaysExistsChecker{}, publisher)

	if _, err := uc.EnrollCourse(context.Background(), "user-1", "course-1"); err != nil {
		t.Fatalf("enroll failed: %v", err)
	}
	progress, err := uc.MarkLessonCompleted(context.Background(), "user-1", "course-1", "lesson-1", 1)
	if err != nil {
		t.Fatalf("mark completed failed: %v", err)
	}
	if progress.Percentage != 100 {
		t.Fatalf("expected 100 percent progress, got %v", progress.Percentage)
	}

	found := false
	for _, subject := range publisher.subjects {
		if subject == EventCourseCompleted {
			found = true
		}
	}
	if !found {
		t.Fatal("expected course.completed event")
	}
}
