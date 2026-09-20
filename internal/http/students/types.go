package students

import "github.com/google/uuid"

type StudentDetails struct {
	ID    uuid.UUID `json:"_id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
