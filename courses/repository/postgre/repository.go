package postgre

import (
	"context"
	"database/sql"
	"errors"

	"online-learning-platform/internal/services/courses/model"
	"online-learning-platform/internal/services/courses/repository"
)

type Repository struct{ db *sql.DB }

func New(db *sql.DB) repository.CourseRepository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, c model.Course) (model.Course, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO courses (id,title,description,instructor_id,published,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`, c.ID, c.Title, c.Description, c.InstructorID, c.Published, c.CreatedAt, c.UpdatedAt)
	return c, err
}

func (r *Repository) GetByID(ctx context.Context, id string) (model.Course, error) {
	var c model.Course
	err := r.db.QueryRowContext(ctx, `SELECT id,title,description,instructor_id,published,created_at,updated_at FROM courses WHERE id=$1`, id).Scan(&c.ID, &c.Title, &c.Description, &c.InstructorID, &c.Published, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Course{}, errors.New("course not found")
	}
	return c, err
}

func (r *Repository) Update(ctx context.Context, c model.Course) (model.Course, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE courses SET title=$1,description=$2,published=$3,updated_at=$4 WHERE id=$5`, c.Title, c.Description, c.Published, c.UpdatedAt, c.ID)
	return c, err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM courses WHERE id=$1`, id)
	return err
}
func (r *Repository) List(ctx context.Context) ([]model.Course, error) { return r.Search(ctx, "") }
func (r *Repository) Search(ctx context.Context, q string) ([]model.Course, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,title,description,instructor_id,published,created_at,updated_at FROM courses WHERE title ILIKE '%' || $1 || '%'`, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Course, 0)
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.InstructorID, &c.Published, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}
func (r *Repository) CreateModule(ctx context.Context, m model.Module) (model.Module, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO modules (id,course_id,title,position) VALUES ($1,$2,$3,$4)`, m.ID, m.CourseID, m.Title, m.Position)
	return m, err
}
func (r *Repository) UpdateModule(ctx context.Context, m model.Module) (model.Module, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE modules SET title=$1,position=$2 WHERE id=$3`, m.Title, m.Position, m.ID)
	return m, err
}
func (r *Repository) DeleteModule(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM modules WHERE id=$1`, id)
	return err
}
func (r *Repository) CreateLesson(ctx context.Context, l model.Lesson) (model.Lesson, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO lessons (id,course_id,module_id,title,content,position) VALUES ($1,$2,$3,$4,$5,$6)`, l.ID, l.CourseID, l.ModuleID, l.Title, l.Content, l.Position)
	return l, err
}
func (r *Repository) UpdateLesson(ctx context.Context, l model.Lesson) (model.Lesson, error) {
	_, err := r.db.ExecContext(ctx, `UPDATE lessons SET title=$1,content=$2,position=$3 WHERE id=$4`, l.Title, l.Content, l.Position, l.ID)
	return l, err
}
func (r *Repository) DeleteLesson(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM lessons WHERE id=$1`, id)
	return err
}
func (r *Repository) GetLesson(ctx context.Context, id string) (model.Lesson, error) {
	var l model.Lesson
	err := r.db.QueryRowContext(ctx, `SELECT id,course_id,module_id,title,content,position FROM lessons WHERE id=$1`, id).Scan(&l.ID, &l.CourseID, &l.ModuleID, &l.Title, &l.Content, &l.Position)
	return l, err
}
