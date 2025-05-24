package dto

// PasswordResetRequest represents the request to initiate password reset
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents the request to confirm password reset
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// EmailVerificationRequest represents the request to verify email
type EmailVerificationRequest struct {
	Token string `json:"token" binding:"required"`
}
