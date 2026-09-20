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

func (h *Handler) AddStudent(c *gin.Context) {
	var req AddStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ApiError(
			c,
			http.StatusBadRequest,
			"Invalid request schema",
		)
		return
	}

	classId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.ApiError(c, http.StatusNotFound, "Class not found")
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

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		response.ApiError(c, http.StatusNotFound, "Student not found")
		return
	}

	class, err := h.service.AddStudent(c.Request.Context(), classId, teacherID, studentID)

	if err != nil {
		h.handleError(c, err)
		return
	}

	response.ApiResponse(
		c,
		http.StatusOK,
		class,
	)
}

func (h *Handler) GetClassDetails(c *gin.Context) {
	classId, err := uuid.Parse(c.Param("id"))

	if err != nil {
		response.ApiError(c, http.StatusNotFound, "Class not found")
		return
	}

	userID, ok := c.MustGet("userId").(uuid.UUID)
	if !ok {
		response.ApiError(c, http.StatusUnauthorized, "Unauthorized, token missing or invalid")
		return
	}

	class, err := h.service.GetClassDetails(c.Request.Context(), classId, userID)

	if err != nil {
		h.handleError(c, err)
		return
	}

	response.ApiResponse(c, http.StatusOK, class)
}

func (h *Handler) MyAttendance(c *gin.Context) {
	classID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.ApiError(c, http.StatusNotFound, "Class not found")
		return
	}

	studentID , ok := c.MustGet("userId").(uuid.UUID)

	if !ok {
			response.ApiError(c, http.StatusUnauthorized, "Unauthorized, token missing or invalid")
			return
		}
	
	result, err := h.service.GetMyAttendance(c.Request.Context(), classID, studentID)

	if err != nil {
		h.handleError(c, err)
		return
	}
	
	response.ApiResponse(c, http.StatusOK, result)
}
