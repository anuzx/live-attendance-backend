package class

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateClass(

	ctx context.Context,
	className string,
	teacherID uuid.UUID,

) (*Class, error) {

	var class Class

	query := `
    INSERT INTO classes (class_name, teacher_id)
    VALUES ($1 , $2)
    RETURNING id, class_name, teacher_id
    `

	err := r.db.QueryRow(
		ctx,
		query,
		className,
		teacherID,
	).Scan(
		&class.ID,
		&class.ClassName,
		&class.TeacherID,
	)

	if err != nil {
		return nil, err
	}

	//new class has no students
	class.StudentIDs = []uuid.UUID{}

	return &class, nil
}
