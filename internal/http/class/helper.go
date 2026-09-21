package class

import (
	"errors"
	"log"
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrClassNotFound):
		response.ApiError(c, http.StatusNotFound, "Class not found")
	case errors.Is(err, ErrStudentNotFound):
		response.ApiError(c, http.StatusNotFound, "Student not found")
	case errors.Is(err, ErrNotClassTeacher):
		response.ApiError(c, http.StatusForbidden, "Forbidden, not class teacher")
	case errors.Is(err, ErrNotEnrolled):
		response.ApiError(c, http.StatusForbidden, "Forbidden, not enrolled in class")
	default:
		log.Printf("internal error: %v", err) // log the real cause
		response.ApiError(c, http.StatusInternalServerError, "Internal server error")
	}
}
