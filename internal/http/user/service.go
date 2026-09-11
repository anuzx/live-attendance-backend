package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Signup(
	ctx context.Context,
	name string,
	email string,
	password string,
	role string,
) (*User, error) {

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	return s.repository.CreateUser(
		ctx,
		name,
		email,
		string(passwordHash),
		role,
	)
}