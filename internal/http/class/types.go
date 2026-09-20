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

type StudentDetails struct {
	ID    uuid.UUID `json:"_id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type ClassDetails struct {
	ID        uuid.UUID        `json:"_id"`
	ClassName string           `json:"className"`
	TeacherID uuid.UUID        `json:"teacherId"`
	Students  []StudentDetails `json:"students"`
}

type MyAttendance struct {
	ClassID uuid.UUID `json:"classId"`
	Status  *string   `json:"status"`
}
