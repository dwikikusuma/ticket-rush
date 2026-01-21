package handler

import (
	"net/http"

	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/domain"
	"github.com/dwikikusuma/ticket-rush/services/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/login", h.login)
	r.POST("/register", h.register)
}

func (h *AuthHandler) login(c *gin.Context) {
	var req domain.LoginRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	token, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) register(c *gin.Context) {
	var req domain.LoginRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	err := h.svc.Register(ctx, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}
