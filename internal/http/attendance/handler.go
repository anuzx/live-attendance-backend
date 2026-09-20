package attendance

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
	return &Handler{service: service}
}

func (h *Handler) StartSession(c *gin.Context) {
	var req StartSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ApiError(c, http.StatusBadRequest, "Invalid request schema")
		return
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, "Class not found")
		return
	}

	teacherID, ok := c.MustGet("userId").(uuid.UUID)
	if !ok {
		response.ApiError(c, http.StatusUnauthorized, "Unauthorized, token missing or invalid")
		return
	}

	result, err := h.service.StartSession(c.Request.Context(), classID, teacherID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.ApiResponse(c, http.StatusOK, result)
}
