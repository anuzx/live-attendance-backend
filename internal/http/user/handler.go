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

	token, err := h.service.Login(
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
		},
	)

}

func (h *Handler) Me(c *gin.Context) {

	value, exists := c.Get("userId")

	if !exists {
		response.ApiError(
			c,
			http.StatusUnauthorized,
			"Unauthorized, token missing or invalid",
		)
		return
	}

	userID, ok := value.(uuid.UUID)

	if !ok {
		response.ApiError(
			c,
			http.StatusUnauthorized,
			"Unauthorized, token missing or invalid",
		)
		return
	}

	user, err := h.service.GetMe(c.Request.Context(), userID)

	if err != nil {
		response.ApiError(c, http.StatusUnauthorized, "Unauthorized, token missing or invalid")

		return
	}

	response.ApiResponse(
		c,
		http.StatusOK,
		SignupResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	)
}
