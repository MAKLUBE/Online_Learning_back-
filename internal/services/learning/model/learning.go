package model

import "time"

type Enrollment struct {
	ID        string
	UserID    string
	CourseID  string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Progress struct {
	UserID           string
	CourseID         string
	CompletedLessons int
	TotalLessons     int
	Percentage       float64
	UpdatedAt        time.Time
}

type Assignment struct {
	ID       string
	CourseID string
	LessonID string
	Title    string
	DueAt    time.Time
}

type Submission struct {
	ID           string
	AssignmentID string
	UserID       string
	Answer       string
	Grade        float64
	Status       string
	SubmittedAt  time.Time
	GradedAt     time.Time
}
