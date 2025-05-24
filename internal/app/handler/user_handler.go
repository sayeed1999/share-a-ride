package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sayeed1999/share-a-ride/internal/app/dto"
	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/pkg/dateutil"
	"github.com/sayeed1999/share-a-ride/internal/pkg/hashutil"
	"github.com/sayeed1999/share-a-ride/internal/pkg/jwtutil"
	"github.com/sayeed1999/share-a-ride/internal/provider/email"
)

type UserHandler struct {
	userUseCase  usecase.UserUseCase
	emailService email.EmailServiceInterface
}

func NewUserHandler(userUseCase usecase.UserUseCase, emailService email.EmailServiceInterface) *UserHandler {
	return &UserHandler{
		userUseCase:  userUseCase,
		emailService: emailService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := hashutil.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	// Generate verification token
	verifyToken := uuid.New().String()

	user := &entity.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashedPassword,
		Role:        entity.RoleUser,
		VerifyToken: verifyToken,
	}

	if err := h.userUseCase.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send verification email
	if err := h.emailService.SendVerificationEmail(user.Email, verifyToken); err != nil {
		// Log error but don't return it to user
		// Consider implementing retry mechanism
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: dateutil.FormatDateTime(user.CreatedAt),
		UpdatedAt: dateutil.FormatDateTime(user.UpdatedAt),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token is required"})
		return
	}

	user, err := h.userUseCase.GetUserByVerifyToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invalid verification token"})
		return
	}

	user.IsEmailVerified = true
	user.VerifyToken = "" // Clear the token

	if err := h.userUseCase.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified successfully"})
}

func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req dto.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists
		c.JSON(http.StatusOK, gin.H{"message": "if your email is registered, you will receive a password reset link"})
		return
	}

	// Generate reset token
	resetToken := uuid.New().String()
	user.ResetToken = resetToken
	user.ResetExpires = time.Now().Add(1 * time.Hour)

	if err := h.userUseCase.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process request"})
		return
	}

	// Send reset email
	if err := h.emailService.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		// Log error but don't return it to user
	}

	c.JSON(http.StatusOK, gin.H{"message": "if your email is registered, you will receive a password reset link"})
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req dto.PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.GetUserByResetToken(req.Token)
	if err != nil || time.Now().After(user.ResetExpires) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired reset token"})
		return
	}

	// Hash new password
	hashedPassword, err := hashutil.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user.Password = hashedPassword
	user.ResetToken = ""
	user.ResetExpires = time.Time{}

	if err := h.userUseCase.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.userUseCase.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: dateutil.FormatDateTime(user.CreatedAt),
		UpdatedAt: dateutil.FormatDateTime(user.UpdatedAt),
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Check password
	if !hashutil.CheckPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := jwtutil.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	response := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: dateutil.FormatDateTime(user.CreatedAt),
			UpdatedAt: dateutil.FormatDateTime(user.UpdatedAt),
		},
	}

	c.JSON(http.StatusOK, response)
}
