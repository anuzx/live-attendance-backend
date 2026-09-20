package auth

import (
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, role, err := ParseToken(c.GetHeader("Authorization"))
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
		c.Set("role", role)

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
