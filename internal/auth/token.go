package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userId uuid.UUID, role string) (string, error) {

	claims := jwt.MapClaims{
		"userId": userId.String(),
		"role":   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secretkey"))
}
