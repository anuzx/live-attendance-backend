package students

import (
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}
func (h *Handler) GetStudents(c *gin.Context) {

	students, err := h.service.GetAllStudents(c.Request.Context())
	if err != nil {
		return
	}

	response.ApiResponse(c, http.StatusOK, students)
}
