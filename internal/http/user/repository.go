package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailExists = errors.New("email already exists")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(
	ctx context.Context,
	name string,
	email string,
	password string,
	role string,
) (*User, error) {

	query := `
      INSERT INTO users(
         name,
         email,
         password,
         role
      )
      VALUES ($1 , $2 , $3 , $4)
      RETURNING id, name, email, role
  `

	var user User

	err := r.db.QueryRow(
		ctx,
		query,
		name,
		email,
		password,
		role,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" {
			return nil, ErrEmailExists
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("failed to create user")
		}

		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserBYEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	query :=
		`
	SELECT id,name,email,password,role
	FROM users
	WHERE email =$1
	`

	var user User

	//.scan() -> standard method used to copy data from database rows into go variables

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return &user, nil

}
