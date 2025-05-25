package config

// TODO: Implement tests for the configuration loading and validation

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: - [BUG] Tests are not working, config.Load() is not getting the mock env values.
func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("APP_BASE_URL", "http://test-server:8080")
	os.Setenv("APP_NAME", "ShareARide")
	os.Setenv("APP_VERSION", "1.0.0")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "test-server")
	os.Setenv("DATABASE_HOST", "test-server")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("DATABASE_USER", "testuser")
	os.Setenv("DATABASE_PASSWORD", "testpassword")
	os.Setenv("DATABASE_NAME", "shareride")
	os.Setenv("DATABASE_SSLMODE", "disable")

	// Load the configuration
	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Validate the loaded configuration
	assert.Equal(t, "http://test-server:8080", cfg.App.BaseURL)
	assert.Equal(t, "ShareARide", cfg.App.Name)
	assert.Equal(t, "1.0.0", cfg.App.Version)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "test-server", cfg.Server.Host)
	assert.Equal(t, "test-server", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpassword", cfg.Database.Password)
	assert.Equal(t, "shareride", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
}
