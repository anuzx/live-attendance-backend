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

func (s *Service) AddStudent(
	ctx context.Context,
	classID, teacherID, studentID uuid.UUID,
) (*Class, error) {

	class, err := s.repository.GetClassByID(ctx, classID)
	if err != nil {
		return nil, ErrClassNotFound
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

func (s *Service) GetClassDetails(
	ctx context.Context,
	classID uuid.UUID,
	userID uuid.UUID,
) (*ClassDetails, error) {

	class, err := s.repository.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}

	students, err := s.repository.GetStudentByClassID(ctx, classID)
	if err != nil {
		return nil, err
	}
	allowed := class.TeacherID == userID
	if !allowed {
		for _, st := range students {
			if st.ID == userID {
				allowed = true
				break
			}
		}
	}

	if !allowed {
		return nil, ErrNotClassTeacher
	}

	return &ClassDetails{
		ID:        class.ID,
		ClassName: class.ClassName,
		TeacherID: class.TeacherID,
		Students:  students,
	}, nil
}

func (s *Service) GetMyAttendance(
	ctx context.Context,
	classID, studentID uuid.UUID,
) (*MyAttendance, error) {

	if _, err := s.repository.GetClassByID(ctx, classID); err != nil {
		return nil, err
	}

	enrolled, err := s.repository.IsStudentEnrolled(ctx, classID, studentID)

	if err != nil {
		return nil, err
	}
	if !enrolled {
		return nil, err
	}

	status, err := s.repository.GetAttendanceStatus(ctx, classID, studentID)
	if err != nil {
		return nil, err
	}

	return &MyAttendance{ClassID: classID, Status: status}, nil
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