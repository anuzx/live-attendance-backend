package attendance

import (
	"context"

	"github.com/anuzx/live-attendance-backend/internal/http/class"
	"github.com/google/uuid"
)

type Service struct {
	classService *class.Service
	store        *SessionStore
}

func NewService(classService *class.Service, store *SessionStore) *Service {
	return &Service{classService: classService, store: store}
}

func (s *Service) StartSession(
	ctx context.Context,
	classID, teacherID uuid.UUID,
) (*StartSessionResponse, error) {

	c, err := s.classService.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err // class.ErrClassNotFound or a DB error
	}

	if c.TeacherID != teacherID {
		return nil, class.ErrNotClassTeacher
	}

	session := s.store.Start(classID)

	return &StartSessionResponse{
		ClassID:   session.ClassID.String(),
		StartedAt: session.StartedAt,
	}, nil
}
