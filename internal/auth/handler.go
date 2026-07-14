package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Validation failed",
		})
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"code":    "EMAIL_EXISTS",
				"message": "An account with this email already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to register user",
		})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_INPUT",
			"message": "Validation failed",
		})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    "INVALID_CREDENTIALS",
				"message": "Email or password is incorrect",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "Failed to login",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
