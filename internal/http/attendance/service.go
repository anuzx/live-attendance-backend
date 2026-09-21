package attendance

import (
	"context"

	"github.com/anuzx/live-attendance-backend/internal/http/class"
	"github.com/google/uuid"
)

type Service struct {
	classService *class.Service
	store        *SessionStore
	repository   *Repository
}

func NewService(classService *class.Service, store *SessionStore, repository *Repository) *Service {
	return &Service{
		classService: classService,
		store:        store,
		repository: repository,
	}
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

type FinalSummary struct {
	Present int
	Absent  int
	Total   int
}

func (s *Service) FinishSession(ctx context.Context) (*FinalSummary, error) {

	// snapshot of the active session
	session, ok := s.store.Get()
	if !ok {
		return nil, ErrNoActiveSession
	}

	// step 2: everyone enrolled in the class
	enrolled, err := s.repository.GetEnrolledStudentIDs(ctx, session.ClassID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(enrolled))
	statuses := make([]string, 0, len(enrolled))
	summary := &FinalSummary{}

	for _, id := range enrolled {
		key := id.String()

		status, marked := session.Attendance[key]
		if !marked {
			status = "absent" // step 3: unmarked students become absent
		}

		if status == "present" {
			summary.Present++ // step 5: final summary
		} else {
			summary.Absent++
		}

		ids = append(ids, key)
		statuses = append(statuses, status)
	}
	summary.Total = len(enrolled)

	// step 4: persist
	if err := s.repository.SaveAttendance(ctx, session.ClassID, ids, statuses); err != nil {
		return nil, err
	}

	// step 6: clear memory, only after the save succeeded
	s.store.Clear()

	return summary, nil
}

// IsStudentEnrolled reports whether the student is enrolled in the class.
// Used by the websocket layer to validate ATTENDANCE_MARKED events.
func (s *Service) IsStudentEnrolled(
	ctx context.Context,
	classID, studentID uuid.UUID,
) (bool, error) {

	ids, err := s.repository.GetEnrolledStudentIDs(ctx, classID)
	if err != nil {
		return false, err
	}

	for _, id := range ids {
		if id == studentID {
			return true, nil
		}
	}

	return false, nil
}
