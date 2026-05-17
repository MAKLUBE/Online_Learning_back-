package repository

import (
	"context"

	"online-learning-platform/internal/services/courses/model"
)

type CourseRepository interface {
	Create(ctx context.Context, course model.Course) (model.Course, error)
	GetByID(ctx context.Context, id string) (model.Course, error)
	Update(ctx context.Context, course model.Course) (model.Course, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Course, error)
	Search(ctx context.Context, query string) ([]model.Course, error)
	CreateModule(ctx context.Context, module model.Module) (model.Module, error)
	UpdateModule(ctx context.Context, module model.Module) (model.Module, error)
	DeleteModule(ctx context.Context, id string) error
	CreateLesson(ctx context.Context, lesson model.Lesson) (model.Lesson, error)
	UpdateLesson(ctx context.Context, lesson model.Lesson) (model.Lesson, error)
	DeleteLesson(ctx context.Context, id string) error
	GetLesson(ctx context.Context, id string) (model.Lesson, error)
}
