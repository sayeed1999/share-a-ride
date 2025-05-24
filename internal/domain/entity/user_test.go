package entity

import (
	"testing"
	"time"
)

func TestUserRoles(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected bool
	}{
		{
			name:     "Valid admin role",
			role:     RoleAdmin,
			expected: true,
		},
		{
			name:     "Valid user role",
			role:     RoleUser,
			expected: true,
		},
		{
			name:     "Valid guest role",
			role:     RoleGuest,
			expected: true,
		},
		{
			name:     "Invalid role",
			role:     Role("invalid"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidRole(tt.role)
			if isValid != tt.expected {
				t.Errorf("isValidRole(%s) = %v, want %v", tt.role, isValid, tt.expected)
			}
		})
	}
}

func TestUserCreation(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "Valid user",
			user: User{
				Email:           "test@example.com",
				Password:        "password123",
				Name:            "Test User",
				Role:            RoleUser,
				IsEmailVerified: false,
				CreatedAt:       now,
				UpdatedAt:       now,
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			user: User{
				Email:    "invalid-email",
				Password: "password123",
				Name:     "Test User",
				Role:     RoleUser,
			},
			wantErr: true,
		},
		{
			name: "Empty password",
			user: User{
				Email: "test@example.com",
				Name:  "Test User",
				Role:  RoleUser,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUser(&tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function to validate roles
func isValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

// Helper function to validate user
func validateUser(u *User) error {
	// Add validation logic here
	return nil
}
