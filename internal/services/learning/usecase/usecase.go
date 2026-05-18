package usecase

import (
	"context"
	"errors"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/learning/model"
	"online-learning-platform/internal/services/learning/repository"
)

type UserChecker interface {
	Exists(ctx context.Context, userID string) (bool, error)
}

type CourseChecker interface {
	Exists(ctx context.Context, courseID string) (bool, error)
}

type UseCase struct {
	repo      repository.Repository
	users     UserChecker
	courses   CourseChecker
	publisher events.Publisher
}

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrCourseNotFound = errors.New("course not found")
)

func New(repo repository.Repository, users UserChecker, courses CourseChecker, publisher events.Publisher) *UseCase {
	if publisher == nil {
		publisher = events.NoopPublisher{}
	}
	return &UseCase{repo: repo, users: users, courses: courses, publisher: publisher}
}

func (u *UseCase) EnrollCourse(ctx context.Context, userID, courseID string) (model.Enrollment, error) {
	if u.users != nil {
		ok, err := u.users.Exists(ctx, userID)
		if err != nil {
			return model.Enrollment{}, err
		}
		if !ok {
			return model.Enrollment{}, ErrUserNotFound
		}
	}
	if u.courses != nil {
		ok, err := u.courses.Exists(ctx, courseID)
		if err != nil {
			return model.Enrollment{}, err
		}
		if !ok {
			return model.Enrollment{}, ErrCourseNotFound
		}
	}
	now := time.Now().UTC()
	enrollment, err := u.repo.CreateEnrollment(ctx, model.Enrollment{
		ID:        events.NewID(),
		UserID:    userID,
		CourseID:  courseID,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return model.Enrollment{}, err
	}
	_, _ = u.repo.SaveProgress(ctx, model.Progress{UserID: userID, CourseID: courseID, TotalLessons: 1, UpdatedAt: now})
	_ = u.publisher.Publish(ctx, EventCourseEnrolled, events.New(EventCourseEnrolled, map[string]any{
		"user_id":   userID,
		"course_id": courseID,
	}))
	return enrollment, nil
}

func (u *UseCase) UnenrollCourse(ctx context.Context, userID, courseID string) error {
	return u.repo.DeleteEnrollment(ctx, userID, courseID)
}

func (u *UseCase) ListMyCourses(ctx context.Context, userID string) ([]model.Enrollment, error) {
	return u.repo.ListEnrollmentsByUser(ctx, userID)
}

func (u *UseCase) GetEnrollment(ctx context.Context, userID, courseID string) (model.Enrollment, error) {
	return u.repo.GetEnrollment(ctx, userID, courseID)
}

func (u *UseCase) MarkLessonCompleted(ctx context.Context, userID, courseID, lessonID string, totalLessons int) (model.Progress, error) {
	progress, err := u.repo.GetProgress(ctx, userID, courseID)
	if err != nil {
		progress = model.Progress{UserID: userID, CourseID: courseID}
	}
	if totalLessons <= 0 {
		totalLessons = progress.TotalLessons
	}
	if totalLessons <= 0 {
		totalLessons = 1
	}
	progress.TotalLessons = totalLessons
	progress.CompletedLessons++
	if progress.CompletedLessons > progress.TotalLessons {
		progress.CompletedLessons = progress.TotalLessons
	}
	progress.Percentage = float64(progress.CompletedLessons) / float64(progress.TotalLessons) * 100
	progress.UpdatedAt = time.Now().UTC()
	saved, err := u.repo.SaveProgress(ctx, progress)
	if err != nil {
		return model.Progress{}, err
	}
	_ = u.publisher.Publish(ctx, EventLessonCompleted, events.New(EventLessonCompleted, map[string]any{
		"user_id": userID, "course_id": courseID, "lesson_id": lessonID, "percentage": saved.Percentage,
	}))
	if saved.Percentage >= 100 {
		_ = u.publisher.Publish(ctx, EventCourseCompleted, events.New(EventCourseCompleted, map[string]any{
			"user_id": userID, "course_id": courseID,
		}))
	}
	return saved, nil
}

func (u *UseCase) GetCourseProgress(ctx context.Context, userID, courseID string) (model.Progress, error) {
	return u.repo.GetProgress(ctx, userID, courseID)
}

func (u *UseCase) GetUserProgress(ctx context.Context, userID string) ([]model.Progress, error) {
	return u.repo.ListProgressByUser(ctx, userID)
}

func (u *UseCase) CreateAssignment(ctx context.Context, courseID, lessonID, title string, dueAt time.Time) (model.Assignment, error) {
	return u.repo.CreateAssignment(ctx, model.Assignment{ID: events.NewID(), CourseID: courseID, LessonID: lessonID, Title: title, DueAt: dueAt})
}

func (u *UseCase) UpdateAssignment(ctx context.Context, assignment model.Assignment) (model.Assignment, error) {
	return u.repo.UpdateAssignment(ctx, assignment)
}

func (u *UseCase) SubmitAssignment(ctx context.Context, assignmentID, userID, answer string) (model.Submission, error) {
	submission, err := u.repo.CreateSubmission(ctx, model.Submission{
		ID: events.NewID(), AssignmentID: assignmentID, UserID: userID, Answer: answer, Status: "submitted", SubmittedAt: time.Now().UTC(),
	})
	if err != nil {
		return model.Submission{}, err
	}
	_ = u.publisher.Publish(ctx, EventAssignmentSubmitted, events.New(EventAssignmentSubmitted, map[string]any{
		"submission_id": submission.ID, "assignment_id": assignmentID, "user_id": userID,
	}))
	return submission, nil
}

func (u *UseCase) GetSubmission(ctx context.Context, id string) (model.Submission, error) {
	return u.repo.GetSubmission(ctx, id)
}

func (u *UseCase) ListSubmissions(ctx context.Context, assignmentID string) ([]model.Submission, error) {
	return u.repo.ListSubmissions(ctx, assignmentID)
}

func (u *UseCase) GradeAssignment(ctx context.Context, submissionID string, grade float64) (model.Submission, error) {
	submission, err := u.repo.GetSubmission(ctx, submissionID)
	if err != nil {
		return model.Submission{}, err
	}
	submission.Grade = grade
	submission.Status = "graded"
	submission.GradedAt = time.Now().UTC()
	updated, err := u.repo.UpdateSubmission(ctx, submission)
	if err != nil {
		return model.Submission{}, err
	}
	_ = u.publisher.Publish(ctx, EventAssignmentGraded, events.New(EventAssignmentGraded, map[string]any{
		"submission_id": updated.ID, "assignment_id": updated.AssignmentID, "user_id": updated.UserID, "grade": updated.Grade,
	}))
	return updated, nil
}

func (u *UseCase) GetGrades(ctx context.Context, userID string) ([]model.Submission, error) {
	return u.repo.ListGrades(ctx, userID)
}
