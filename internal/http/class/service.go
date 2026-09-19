package class

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateClass(
	ctx context.Context,
	className string,
	teacherID uuid.UUID,
) (*Class, error) {

	return s.repository.CreateClass(
		ctx,
		className,
		teacherID,
	)
}
