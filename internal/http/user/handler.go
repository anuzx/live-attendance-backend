package user

import "github.com/google/uuid"

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required , min=6"`
	Role     string `json:"role" binding:"required,oneof=student teacher"`
}

type SignupResponse struct {
	ID    uuid.UUID `json:"_id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
