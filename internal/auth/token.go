package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("secretkey") 

func GenerateToken(userId uuid.UUID, role string) (string, error) {

	claims := jwt.MapClaims{
		"userId": userId.String(),
		"role":   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}


func ParseToken(tokenString string) (uuid.UUID, string, error) {

	//jwt.parse -> splits the token into header, payload, signature and decodes the header and payload 
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) { return jwtSecret, nil },
	)
	if err != nil || !token.Valid {
		return uuid.Nil, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, "", errors.New("invalid claims")
	}

	idStr, ok := claims["userId"].(string)
	role, ok2 := claims["role"].(string)
	if !ok || !ok2 {
		return uuid.Nil, "", errors.New("invalid claims")
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, "", errors.New("invalid claims")
	}

	return userID, role, nil
}