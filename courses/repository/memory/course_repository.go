package memory

import (
	"context"
	"errors"
	"strings"
	"sync"

	"online-learning-platform/internal/services/courses/model"
)

type CourseRepository struct {
	mu      sync.RWMutex
	courses map[string]model.Course
	modules map[string]model.Module
	lessons map[string]model.Lesson
}

func NewCourseRepository() *CourseRepository {
	return &CourseRepository{
		courses: map[string]model.Course{
			"course-1": {ID: "course-1", Title: "Template Course", Description: "Replace this with real storage and logic.", Published: true},
		},
		modules: make(map[string]model.Module),
		lessons: make(map[string]model.Lesson),
	}
}

func (r *CourseRepository) Create(_ context.Context, course model.Course) (model.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.courses[course.ID] = course
	return course, nil
}

func (r *CourseRepository) GetByID(_ context.Context, id string) (model.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	course, ok := r.courses[id]
	if !ok {
		return model.Course{}, errors.New("course not found")
	}
	return course, nil
}

func (r *CourseRepository) Update(_ context.Context, course model.Course) (model.Course, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.courses[course.ID]; !ok {
		return model.Course{}, errors.New("course not found")
	}
	r.courses[course.ID] = course
	return course, nil
}

func (r *CourseRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.courses, id)
	return nil
}

func (r *CourseRepository) List(context.Context) ([]model.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	courses := make([]model.Course, 0, len(r.courses))
	for _, course := range r.courses {
		courses = append(courses, course)
	}
	return courses, nil
}

func (r *CourseRepository) Search(_ context.Context, query string) ([]model.Course, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	query = strings.ToLower(query)
	courses := make([]model.Course, 0)
	for _, course := range r.courses {
		if query == "" || strings.Contains(strings.ToLower(course.Title), query) {
			courses = append(courses, course)
		}
	}
	return courses, nil
}

func (r *CourseRepository) CreateModule(_ context.Context, module model.Module) (model.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules[module.ID] = module
	return module, nil
}

func (r *CourseRepository) UpdateModule(_ context.Context, module model.Module) (model.Module, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.modules[module.ID] = module
	return module, nil
}

func (r *CourseRepository) DeleteModule(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.modules, id)
	return nil
}

func (r *CourseRepository) CreateLesson(_ context.Context, lesson model.Lesson) (model.Lesson, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lessons[lesson.ID] = lesson
	return lesson, nil
}

func (r *CourseRepository) UpdateLesson(_ context.Context, lesson model.Lesson) (model.Lesson, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lessons[lesson.ID] = lesson
	return lesson, nil
}

func (r *CourseRepository) DeleteLesson(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.lessons, id)
	return nil
}

func (r *CourseRepository) GetLesson(_ context.Context, id string) (model.Lesson, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lesson, ok := r.lessons[id]
	if !ok {
		return model.Lesson{}, errors.New("lesson not found")
	}
	return lesson, nil
}
