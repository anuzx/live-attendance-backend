package auth

import (
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			response.ApiError(
				c,
				http.StatusUnauthorized,
				"Unauthorized, token missing or invalid",
			)
			c.Abort()
			return
		}

		decode, err := jwt.Parse(
			token,
			func(t *jwt.Token) (any, error) {
				return []byte("secretkey"), nil
			},
		)

		if err != nil || !decode.Valid {
			response.ApiError(
				c,
				http.StatusUnauthorized,
				"Unauthorized, token missing or invalid",
			)
			c.Abort()
			return
		}

		claims, ok := decode.Claims.(jwt.MapClaims)

		if !ok {
			response.ApiError(
				c,
				http.StatusUnauthorized,
				"Unauthorized, token missing or invalid",
			)
			c.Abort()
			return
		}

		userIDString, ok := claims["userId"].(string)

		if !ok {
			response.ApiError(
				c,
				http.StatusUnauthorized,
				"Unauthorized, token missing or invalid",
			)
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDString)

		if err != nil {
			response.ApiError(
				c,
				http.StatusUnauthorized,
				"Unauthorized, token missing or invalid",
			)
			c.Abort()
			return
		}

		c.Set("userId", userID)
		c.Set("role", claims["role"])

		c.Next()
	}
}

func OnlyTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || role != "teacher" {
			response.ApiError(
				c,
				http.StatusForbidden,
				"Forbidden, teacher access required",
			)

			c.Abort()
			return
		}

		c.Next()
	}
}

func OnlyStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || role != "student" {
			response.ApiError(c, http.StatusForbidden, "Forbidden, student access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
