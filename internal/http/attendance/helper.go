package attendance

import (
	"errors"
	"log"
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/http/class"
	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, class.ErrClassNotFound):
		response.ApiError(c, http.StatusNotFound, "Class not found")
	case errors.Is(err, class.ErrNotClassTeacher):
		response.ApiError(c, http.StatusForbidden, "Forbidden, not class teacher")
	default:
		log.Printf("internal error: %v", err)
		response.ApiError(c, http.StatusInternalServerError, "Internal server error")
	}
}
