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

func (s *Service) GetClassByID(
	ctx context.Context,
	id uuid.UUID,

) (*Class, error) {

	return s.repository.GetClassByID(
		ctx,
		id,
	)
}

func (s *Service) GetStudentByID(
	ctx context.Context,
	id uuid.UUID,

) (*Student, error) {

	return s.repository.GetStudentByID(
		ctx,
		id,
	)
}

func (s *Service) AddStudent(
	ctx context.Context,
	classID, teacherID, studentID uuid.UUID,
) (*Class, error) {

	class, err := s.repository.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}

	if class.TeacherID != teacherID {
		return nil, ErrNotClassTeacher
	}

	if _, err := s.repository.GetStudentByID(ctx, studentID); err != nil {
		return nil, err
	}

	if err := s.repository.AddStudent(ctx, classID, studentID); err != nil {
		return nil, err
	}

	ids, err := s.repository.GetStudentIds(ctx, classID)
	if err != nil {
		return nil, err
	}
	class.StudentIDs = ids

	return class, nil
}
