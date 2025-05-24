package email

import (
	"testing"

	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestSendVerificationEmail(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL: "http://localhost:8080",
		},
		Email: config.EmailConfig{
			Host:     "localhost",
			Port:     2525, // Test SMTP port
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		},
	}

	// Create the email service
	emailService := NewEmailService(cfg)

	tests := []struct {
		name    string
		to      string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid email parameters",
			to:      "user@example.com",
			token:   "verify-token-123",
			wantErr: false,
		},
		{
			name:    "Invalid email address",
			to:      "invalid-email",
			token:   "verify-token-123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := emailService.SendVerificationEmail(tt.to, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSendPasswordResetEmail(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		App: config.AppConfig{
			BaseURL: "http://localhost:8080",
		},
		Email: config.EmailConfig{
			Host:     "localhost",
			Port:     2525, // Test SMTP port
			Username: "test@example.com",
			Password: "password123",
			From:     "noreply@example.com",
		},
	}

	// Create the email service
	emailService := NewEmailService(cfg)

	tests := []struct {
		name    string
		to      string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid email parameters",
			to:      "user@example.com",
			token:   "reset-token-123",
			wantErr: false,
		},
		{
			name:    "Invalid email address",
			to:      "invalid-email",
			token:   "reset-token-123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := emailService.SendPasswordResetEmail(tt.to, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
