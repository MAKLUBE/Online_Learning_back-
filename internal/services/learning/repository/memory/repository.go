package memory

import (
	"context"
	"errors"
	"sync"

	"online-learning-platform/internal/services/learning/model"
)

type Repository struct {
	mu          sync.RWMutex
	enrollments map[string]model.Enrollment
	progress    map[string]model.Progress
	assignments map[string]model.Assignment
	submissions map[string]model.Submission
}

func New() *Repository {
	return &Repository{
		enrollments: make(map[string]model.Enrollment),
		progress:    make(map[string]model.Progress),
		assignments: make(map[string]model.Assignment),
		submissions: make(map[string]model.Submission),
	}
}

func key(userID, courseID string) string { return userID + ":" + courseID }

func (r *Repository) CreateEnrollment(_ context.Context, enrollment model.Enrollment) (model.Enrollment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.enrollments[key(enrollment.UserID, enrollment.CourseID)] = enrollment
	return enrollment, nil
}

func (r *Repository) DeleteEnrollment(_ context.Context, userID, courseID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.enrollments, key(userID, courseID))
	return nil
}

func (r *Repository) GetEnrollment(_ context.Context, userID, courseID string) (model.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	enrollment, ok := r.enrollments[key(userID, courseID)]
	if !ok {
		return model.Enrollment{}, errors.New("enrollment not found")
	}
	return enrollment, nil
}

func (r *Repository) ListEnrollmentsByUser(_ context.Context, userID string) ([]model.Enrollment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Enrollment, 0)
	for _, enrollment := range r.enrollments {
		if enrollment.UserID == userID {
			items = append(items, enrollment)
		}
	}
	return items, nil
}

func (r *Repository) SaveProgress(_ context.Context, progress model.Progress) (model.Progress, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.progress[key(progress.UserID, progress.CourseID)] = progress
	return progress, nil
}

func (r *Repository) GetProgress(_ context.Context, userID, courseID string) (model.Progress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	progress, ok := r.progress[key(userID, courseID)]
	if !ok {
		return model.Progress{}, errors.New("progress not found")
	}
	return progress, nil
}

func (r *Repository) ListProgressByUser(_ context.Context, userID string) ([]model.Progress, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Progress, 0)
	for _, progress := range r.progress {
		if progress.UserID == userID {
			items = append(items, progress)
		}
	}
	return items, nil
}

func (r *Repository) CreateAssignment(_ context.Context, assignment model.Assignment) (model.Assignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignments[assignment.ID] = assignment
	return assignment, nil
}

func (r *Repository) UpdateAssignment(_ context.Context, assignment model.Assignment) (model.Assignment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignments[assignment.ID] = assignment
	return assignment, nil
}

func (r *Repository) CreateSubmission(_ context.Context, submission model.Submission) (model.Submission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.submissions[submission.ID] = submission
	return submission, nil
}

func (r *Repository) GetSubmission(_ context.Context, id string) (model.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	submission, ok := r.submissions[id]
	if !ok {
		return model.Submission{}, errors.New("submission not found")
	}
	return submission, nil
}

func (r *Repository) ListSubmissions(_ context.Context, assignmentID string) ([]model.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Submission, 0)
	for _, submission := range r.submissions {
		if submission.AssignmentID == assignmentID {
			items = append(items, submission)
		}
	}
	return items, nil
}

func (r *Repository) UpdateSubmission(_ context.Context, submission model.Submission) (model.Submission, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.submissions[submission.ID] = submission
	return submission, nil
}

func (r *Repository) ListGrades(_ context.Context, userID string) ([]model.Submission, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]model.Submission, 0)
	for _, submission := range r.submissions {
		if submission.UserID == userID && submission.Status == "graded" {
			items = append(items, submission)
		}
	}
	return items, nil
}
