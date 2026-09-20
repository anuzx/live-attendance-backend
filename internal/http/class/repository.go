package class

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *Repository) GetClassByID(
	ctx context.Context,
	id uuid.UUID,
) (*Class, error) {

	query := `
	SELECT id, class_name, teacher_id
	FROM classes
	WHERE id = $1
	`

	var class Class

	//we are selecting 3 columns so we need to scan all 3
	// db.QueryRow -> returns exactly one or zero row
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&class.ID,
		&class.ClassName,
		&class.TeacherID,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrClassNotFound
		}

		return nil, err

	}

	return &class, nil

}

func (r *Repository) GetStudentByID(ctx context.Context, studentId uuid.UUID) (*Student, error) {

	var student Student

	query := `SELECT id FROM users WHERE id=$1`

	err := r.db.QueryRow(
		ctx,
		query,
		studentId,
	).Scan(
		&student.ID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStudentNotFound
		}

		return nil, err
	}

	return &student, nil

}

func (r *Repository) AddStudent(
	ctx context.Context,
	classId uuid.UUID,
	studentId uuid.UUID,
) error {

	query := `
	INSERT INTO class_students(class_id , student_id )
    VALUES ($1 , $2)
    ON CONFLICT (class_id, student_id) DO NOTHING
    `

	//execute statements that do not return rows
	_, err := r.db.Exec(
		ctx,
		query,
		classId,
		studentId,
	)

	return err
}

// retrieving the students are inserting them
func (r *Repository) GetStudentIds(
	ctx context.Context,
	classId uuid.UUID,
) ([]uuid.UUID, error) {

	query := `
	SELECT student_id
	FROM class_students
	WHERE class_id = $1
	`

	//db.Query -> returns zero to many rows
	rows, err := r.db.Query(
		ctx,
		query,
		classId,
	)

	if err != nil {
		return nil, err
	}

	// If you don't close the rows object, the database connection stays open and won't return to the connection pool. This will eventually exhaust your database resources
	defer rows.Close()

	studentIDs := []uuid.UUID{}

	for rows.Next() {

		var studentID uuid.UUID

		if err := rows.Scan(&studentID); err != nil {
			return nil, err
		}

		studentIDs = append(studentIDs, studentID)
	}

	return studentIDs, rows.Err()
}

func (r *Repository) GetStudentByClassID(
	ctx context.Context,
	classID uuid.UUID,
) ([]StudentDetails, error) {

	//fetch all student details in one query, rather than one query per student
	query := `
		SELECT u.id, u.name, u.email
		FROM class_students cs
		JOIN users u
		ON cs.student_id = u.id
		WHERE cs.class_id = $1
		`

	rows, err := r.db.Query(
		ctx,
		query,
		classID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	students := []StudentDetails{} // empty slice , so JSON gives []

	for rows.Next() {
		var s StudentDetails
		if err := rows.Scan(&s.ID, &s.Name, &s.Email); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *Repository) IsStudentEnrolled(ctx context.Context, classID, studentID uuid.UUID) (bool, error) {

	query := `SELECT EXISTS (
	SELECT 1 FROM class_students
	WHERE class_id = $1 AND student_id = $2
	)`

	var enrolled bool

	err := r.db.QueryRow(ctx, query, classID, studentID).Scan(&enrolled)

	return enrolled, err
}

func (r *Repository) GetAttendanceStatus(
	ctx context.Context,
	classID, studentID uuid.UUID,
) (*string, error) {

	query := `
	SELECT status
	FROM attendance
	WHERE class_id = $1 AND student_id = $2
	`

	var status string
	err := r.db.QueryRow(ctx, query, classID, studentID).Scan(&status)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not persisted yet: this is a valid result, not an error
		}
		return nil, err
	}

	return &status, nil
}
