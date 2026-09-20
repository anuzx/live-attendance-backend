package students

import (
	"context"

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

func (r *Repository) GetAllStudents(
	ctx context.Context,

) ([]StudentDetails, error) {

	query := `
	SELECT id, name, email
	FROM users
	WHERE role = 'student'
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	students := []StudentDetails{}

	//loop stop when the rows run out ,rows.Next() returns boolean value ,while loop pattern 
	for rows.Next() {
		var s StudentDetails
		if err := rows.Scan(&s.ID, &s.Name, &s.Email); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()

}
