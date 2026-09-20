package students

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}
func (s *Service) GetAllStudents(ctx context.Context) ([]StudentDetails, error) {
	return s.repository.GetAllStudents(ctx)
}
