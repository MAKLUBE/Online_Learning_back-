package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"online-learning-platform/internal/adapters/events"
	"online-learning-platform/internal/services/courses/model"
	"online-learning-platform/internal/services/courses/repository"
)

var ErrTitleRequired = errors.New("title is required")

type InstructorChecker interface {
	CanManageCourses(ctx context.Context, userID string) (bool, error)
}

type CourseUseCase struct {
	repo              repository.CourseRepository
	publisher         events.Publisher
	instructorChecker InstructorChecker
}

func NewCourseUseCase(repo repository.CourseRepository) *CourseUseCase {
	return &CourseUseCase{repo: repo, publisher: events.NoopPublisher{}}
}

func NewCourseUseCaseWithEvents(repo repository.CourseRepository, publisher events.Publisher) *CourseUseCase {
	if publisher == nil {
		publisher = events.NoopPublisher{}
	}
	return &CourseUseCase{repo: repo, publisher: publisher}
}

func NewCourseUseCaseWithDependencies(repo repository.CourseRepository, publisher events.Publisher, checker InstructorChecker) *CourseUseCase {
	if publisher == nil {
		publisher = events.NoopPublisher{}
	}
	return &CourseUseCase{repo: repo, publisher: publisher, instructorChecker: checker}
}

func (u *CourseUseCase) Create(ctx context.Context, title, description, instructorID string) (model.Course, error) {
	if strings.TrimSpace(title) == "" {
		return model.Course{}, ErrTitleRequired
	}
	if u.instructorChecker != nil {
		ok, err := u.instructorChecker.CanManageCourses(ctx, instructorID)
		if err != nil {
			return model.Course{}, err
		}
		if !ok {
			return model.Course{}, errors.New("instructor_id must belong to instructor or admin")
		}
	}
	now := time.Now().UTC()
	return u.repo.Create(ctx, model.Course{
		ID:           events.NewID(),
		Title:        strings.TrimSpace(title),
		Description:  strings.TrimSpace(description),
		InstructorID: instructorID,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (u *CourseUseCase) GetByID(ctx context.Context, id string) (model.Course, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *CourseUseCase) Update(ctx context.Context, course model.Course) (model.Course, error) {
	course.UpdatedAt = time.Now().UTC()
	return u.repo.Update(ctx, course)
}

func (u *CourseUseCase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *CourseUseCase) List(ctx context.Context) ([]model.Course, error) {
	return u.repo.List(ctx)
}

func (u *CourseUseCase) Search(ctx context.Context, query string) ([]model.Course, error) {
	return u.repo.Search(ctx, strings.TrimSpace(query))
}

func (u *CourseUseCase) Publish(ctx context.Context, id string) (model.Course, error) {
	course, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return model.Course{}, err
	}
	course.Published = true
	course.UpdatedAt = time.Now().UTC()
	updated, err := u.repo.Update(ctx, course)
	if err != nil {
		return model.Course{}, err
	}
	_ = u.publisher.Publish(ctx, "course.published", events.New("course.published", map[string]any{
		"course_id": updated.ID,
		"title":     updated.Title,
	}))
	return updated, nil
}

func (u *CourseUseCase) CreateModule(ctx context.Context, courseID, title string, position int) (model.Module, error) {
	return u.repo.CreateModule(ctx, model.Module{ID: events.NewID(), CourseID: courseID, Title: strings.TrimSpace(title), Position: position})
}

func (u *CourseUseCase) UpdateModule(ctx context.Context, module model.Module) (model.Module, error) {
	return u.repo.UpdateModule(ctx, module)
}

func (u *CourseUseCase) DeleteModule(ctx context.Context, id string) error {
	return u.repo.DeleteModule(ctx, id)
}

func (u *CourseUseCase) CreateLesson(ctx context.Context, courseID, moduleID, title, content string, position int) (model.Lesson, error) {
	return u.repo.CreateLesson(ctx, model.Lesson{ID: events.NewID(), CourseID: courseID, ModuleID: moduleID, Title: strings.TrimSpace(title), Content: content, Position: position})
}

func (u *CourseUseCase) UpdateLesson(ctx context.Context, lesson model.Lesson) (model.Lesson, error) {
	return u.repo.UpdateLesson(ctx, lesson)
}

func (u *CourseUseCase) DeleteLesson(ctx context.Context, id string) error {
	return u.repo.DeleteLesson(ctx, id)
}

func (u *CourseUseCase) GetLesson(ctx context.Context, id string) (model.Lesson, error) {
	return u.repo.GetLesson(ctx, id)
}
