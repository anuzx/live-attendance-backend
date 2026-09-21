package attendance

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNoActiveSession = errors.New("no active attendance session")

//	const activeSession = {
//	  classId: "c101", // current active class
//	  startedAt: "2025-03-11T10:00:00.000Z", // ISO string
//	  attendance: {
//	    "s100": "present",
//	    "s101": "absent"
//	    // studentId: status
//	};
type Session struct {
	ClassID    uuid.UUID
	StartedAt  string            //ISO string
	Attendance map[string]string //studentID -> "present" | "absent"
}

// A RWMutex is a reader/writer mutual exclusion lock. The lock can be held by an arbitrary number of readers or a single writer. The zero value for a RWMutex is an unlocked mutex.
type SessionStore struct {
	mu      sync.RWMutex
	session *Session
}

// factory function
func NewSessionStore() *SessionStore {
	return &SessionStore{}
}

func (s *SessionStore) Start(classID uuid.UUID) Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session = &Session{
		ClassID:    classID,
		StartedAt:  time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		Attendance: map[string]string{},
	}

	return *s.session
}

// Get returns a copy so callers can't touch the map without holding the lock
func (s *SessionStore) Get() (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.session == nil {
		return Session{}, false
	}

	attendanceCopy := make(map[string]string, len(s.session.Attendance))
	for k, v := range s.session.Attendance {
		attendanceCopy[k] = v
	}

	copied := *s.session
	copied.Attendance = attendanceCopy
	return copied, true
}

func (s *SessionStore) MarkAttendance(studentID, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.session == nil {
		return ErrNoActiveSession
	}

	s.session.Attendance[studentID] = status
	return nil
}

func (s *SessionStore) Summary() (present, absent, total int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.session == nil {
		return 0, 0, 0, ErrNoActiveSession
	}

	for _, status := range s.session.Attendance {
		switch status {
		case "present":
			present++
		case "absent":
			absent++
		}
	}

	return present, absent, present + absent, nil
}

func (s *SessionStore) GetStatus(studentID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.session == nil {
		return "", ErrNoActiveSession
	}

	status, ok := s.session.Attendance[studentID]
	if !ok {
		return "not yet updated", nil
	}

	return status, nil
}

func (s *SessionStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.session = nil
}
