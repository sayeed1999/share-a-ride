package config

import (
	"fmt"
	"os"
)

// Config holds all configuration of the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Email    EmailConfig
}

type AppConfig struct {
	BaseURL string
	Name    string
	Version string
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

var cfg *Config

// Load returns a Config struct populated with values from environment variables
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	// Note:- godotenv is requiring a .env file to be present in the root directory,
	// which is not ideal for production environments.
	// In prod, docker or other orchestration tools will handle environment variables.
	// if err := godotenv.Load(); err != nil {
	// 	return nil, fmt.Errorf("error loading .env file: %w", err)
	// }

	cfg = &Config{
		App: AppConfig{
			BaseURL: getEnvOrDefault("APP_BASE_URL", "http://localhost:8080"),
			Name:    getEnvOrDefault("APP_NAME", "Share A Ride"),
			Version: getEnvOrDefault("APP_VERSION", "1.0.0"),
		},
		Server: ServerConfig{
			Port: getEnvOrDefault("SERVER_PORT", "8080"),
			Host: getEnvOrDefault("SERVER_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
			Name:     getEnvOrDefault("DB_NAME", "share_a_ride"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		Email: EmailConfig{
			Host:     getEnvOrDefault("EMAIL_HOST", "smtp.gmail.com"),
			Port:     getEnvAsIntOrDefault("EMAIL_PORT", 587),
			Username: getEnvOrDefault("EMAIL_USERNAME", ""),
			Password: getEnvOrDefault("EMAIL_PASSWORD", ""),
			From:     getEnvOrDefault("EMAIL_FROM", "noreply@sharearide.com"),
		},
	}

	return cfg, nil
}

// Returns the current config if it exists, otherwise loads it
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	return Load()
}

// Returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Name,
		c.Port,
		c.SSLMode,
	)
}

// Returns the full server address
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := fmt.Sscanf(value, "%d"); err == nil {
			return intValue
		}
	}
	return defaultValue
}
