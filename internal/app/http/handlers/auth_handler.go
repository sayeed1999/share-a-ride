package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/domain/errors"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"github.com/sayeed1999/share-a-ride/internal/domain/services"
)

type registerRequest struct {
	Name     string          `json:"name" binding:"required,min=2,max=100"`
	Email    string          `json:"email" binding:"required,email"`
	Phone    string          `json:"phone" binding:"required"`
	Password string          `json:"password" binding:"required,min=8"`
	UserType models.UserType `json:"user_type" binding:"required,oneof=rider driver"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.authService.Register(c.Request.Context(), services.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		UserType: req.UserType,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrEmailExists || err == errors.ErrPhoneExists {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"phone":     user.Phone,
				"user_type": user.UserType,
			},
			"tokens": tokens,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, tokens, err := h.authService.Login(c.Request.Context(), services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrInvalidCredentials {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user": gin.H{
				"id":        user.ID,
				"name":      user.Name,
				"email":     user.Email,
				"phone":     user.Phone,
				"user_type": user.UserType,
			},
			"tokens": tokens,
		},
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrInvalidToken {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"tokens": tokens,
		},
	})
}
