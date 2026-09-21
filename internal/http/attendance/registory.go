package attendance

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

func (r *Repository) GetEnrolledStudentIDs(
	ctx context.Context,
	classID uuid.UUID,
) ([]uuid.UUID, error) {

	rows, err := r.db.Query(ctx,
		`SELECT student_id FROM class_students WHERE class_id = $1`, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) SaveAttendance(
	ctx context.Context,
	classID uuid.UUID,
	studentIDs []string,
	statuses []string,
) error {

	_, err := r.db.Exec(ctx, `
		INSERT INTO attendance (class_id, student_id, status)
		SELECT $1::uuid, u.id, u.status
		FROM unnest($2::uuid[], $3::text[]) AS u(id, status)
		ON CONFLICT (class_id, student_id)
		DO UPDATE SET status = EXCLUDED.status
	`, classID, studentIDs, statuses)

	return err
}
