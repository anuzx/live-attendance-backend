package types

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	password string
	Role     string
}
