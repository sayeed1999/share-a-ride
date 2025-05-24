package email

import (
	"fmt"
	"net/smtp"

	"github.com/sayeed1999/share-a-ride/internal/config"
)

type EmailServiceInterface interface {
	SendVerificationEmail(to, token string) error
	SendPasswordResetEmail(to, token string) error
}

type EmailService struct {
	config *config.Config
}

func NewEmailService(cfg *config.Config) EmailServiceInterface {
	return &EmailService{
		config: cfg,
	}
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Email Verification"
	verifyLink := fmt.Sprintf("%s/verify-email?token=%s", s.config.App.BaseURL, token)
	body := fmt.Sprintf("Please click the link below to verify your email:\n%s", verifyLink)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset Request"
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.config.App.BaseURL, token)
	body := fmt.Sprintf("Please click the link below to reset your password:\n%s\nThis link will expire in 1 hour.", resetLink)

	return s.sendEmail(to, subject, body)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth(
		"",
		s.config.Email.Username,
		s.config.Email.Password,
		s.config.Email.Host,
	)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", s.config.Email.From, to, subject, body)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.Email.Host, s.config.Email.Port),
		auth,
		s.config.Email.From,
		[]string{to},
		[]byte(msg),
	)
}
