package db

import (
	"testing"

	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestInitDB(t *testing.T) {
	// Create test config with test database
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "test_user",
			Password: "test_password",
			DBName:   "test_db",
			SSLMode:  "disable",
		},
	}

	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name:    "Successful database connection",
			config:  cfg,
			wantErr: false,
		},
		{
			name: "Invalid database connection",
			config: &config.Config{
				Database: config.DatabaseConfig{
					Host:     "invalid-host",
					Port:     "5432",
					User:     "invalid-user",
					Password: "invalid-password",
					DBName:   "invalid-db",
					SSLMode:  "disable",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := InitDB(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, DB)

			// Test auto-migration
			var user entity.User
			err = DB.First(&user).Error
			assert.Error(t, err) // Should error since table is empty
			assert.Contains(t, err.Error(), "record not found")

			// Test creating a user
			testUser := &entity.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
				Role:     entity.RoleUser,
			}
			err = DB.Create(testUser).Error
			assert.NoError(t, err)
			assert.NotZero(t, testUser.ID)

			// Clean up
			err = DB.Migrator().DropTable(&entity.User{})
			assert.NoError(t, err)
		})
	}
}
