package user

import (
	"errors"
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

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=student teacher"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupResponse struct {
	ID    uuid.UUID `json:"_id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ApiError(
			c,
			http.StatusBadRequest,
			"Invalid request schema",
		)

		return
	}

	user, err := h.service.Signup(
		c.Request.Context(),
		req.Name,
		req.Email,
		req.Password,
		req.Role,
	)

	if err != nil {
		if errors.Is(err, ErrEmailExists) {
			response.ApiError(
				c,
				http.StatusBadRequest,
				"Email already exists",
			)
			return
		}

		response.ApiError(
			c,
			http.StatusInternalServerError,
			"Failed to create user",
		)
		return
	}

	response.ApiResponse(
		c,
		http.StatusCreated,
		SignupResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ApiError(c, http.StatusBadRequest, "Invalid request schema")
		return

	}

	user, token, err := h.service.Login(
		c.Request.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {

		response.ApiError(
			c,
			http.StatusBadRequest,
			"Invalid email or password",
		)

		return
	}

	response.ApiResponse(
		c,
		http.StatusOK,
		gin.H{
			"token": token,
			"user": SignupResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  user.Role,
			},
		},
	)

}

func Me(c *gin.Context) {

}
