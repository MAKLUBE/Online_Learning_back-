package model

import "time"

type Course struct {
	ID           string
	Title        string
	Description  string
	InstructorID string
	Published    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Module struct {
	ID       string
	CourseID string
	Title    string
	Position int
}

type Lesson struct {
	ID       string
	CourseID string
	ModuleID string
	Title    string
	Content  string
	Position int
}
