package class

import "github.com/google/uuid"

type Class struct {
	ID         uuid.UUID   `json:"_id"`
	ClassName  string      `json:"className"`
	TeacherID  uuid.UUID   `json:"teacherId"`
	StudentIDs []uuid.UUID `json:"studentIds"`
}

type CreateClassRequest struct {
	ClassName string `json:"className" binding:"required"`
}

type AddStudentRequest struct {
	StudentID string `json:"studentId" binding:"required"`
}

type Student struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
	Role     string
}
