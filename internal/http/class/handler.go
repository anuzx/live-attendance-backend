package class

import (
	"net/http"

	"github.com/anuzx/live-attendance-backend/internal/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateClass(c *gin.Context) {

	var req CreateClassRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ApiError(
			c,
			http.StatusBadRequest,
			"Invalid request schema",
		)
		return
	}

	value, exists := c.Get("userId")

	if !exists {
		response.ApiError(
			c,
			http.StatusUnauthorized,
			"Unauthorized, token missing or invalid",
		)

		return 
	}

	teacherID, ok := value.(uuid.UUID)

	if !ok {
		return
	}

	class, err := h.service.CreateClass(
		c.Request.Context(),
		req.ClassName,
		teacherID,
	)

	if err != nil {
		return
	}

	response.ApiResponse(
		c,
		http.StatusCreated,
		class,
	)

}
