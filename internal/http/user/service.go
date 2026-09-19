package user

import (
	"context"
	"errors"

	"github.com/anuzx/live-attendance-backend/internal/auth"
	"github.com/google/uuid"
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

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := s.repository.GetUserBYEmail(ctx, email)

	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	//generate the token after password verification
	token, err := auth.GenerateToken(user.ID, user.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetMe(
	ctx context.Context,
	userID uuid.UUID,
) (*User, error) {

	return s.repository.GetUserByID(ctx, userID)
}
