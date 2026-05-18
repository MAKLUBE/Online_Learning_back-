package repository

import (
	"context"

	"online-learning-platform/internal/services/learning/model"
)

type Repository interface {
	CreateEnrollment(ctx context.Context, enrollment model.Enrollment) (model.Enrollment, error)
	DeleteEnrollment(ctx context.Context, userID, courseID string) error
	GetEnrollment(ctx context.Context, userID, courseID string) (model.Enrollment, error)
	ListEnrollmentsByUser(ctx context.Context, userID string) ([]model.Enrollment, error)
	SaveProgress(ctx context.Context, progress model.Progress) (model.Progress, error)
	GetProgress(ctx context.Context, userID, courseID string) (model.Progress, error)
	ListProgressByUser(ctx context.Context, userID string) ([]model.Progress, error)
	CreateAssignment(ctx context.Context, assignment model.Assignment) (model.Assignment, error)
	UpdateAssignment(ctx context.Context, assignment model.Assignment) (model.Assignment, error)
	CreateSubmission(ctx context.Context, submission model.Submission) (model.Submission, error)
	GetSubmission(ctx context.Context, id string) (model.Submission, error)
	ListSubmissions(ctx context.Context, assignmentID string) ([]model.Submission, error)
	UpdateSubmission(ctx context.Context, submission model.Submission) (model.Submission, error)
	ListGrades(ctx context.Context, userID string) ([]model.Submission, error)
}
