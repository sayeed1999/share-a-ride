package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration of the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Host         string
	Port         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type AppConfig struct {
	BaseURL string
	Name    string
	Version string
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

	cfg = &Config{}

	// Server configuration
	cfg.Server = ServerConfig{
		Host:         getEnv("SERVER_HOST", "localhost"),
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("APP_ENV", "development"),
		ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
		WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
	}

	// Database configuration
	cfg.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "share_a_ride"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// JWT configuration
	cfg.JWT = JWTConfig{
		SecretKey:          getEnv("JWT_SECRET_KEY", "your-256-bit-secret"),
		AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 1*time.Hour),
		RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
	}

	// App configuration
	cfg.App = AppConfig{
		BaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		Name:    getEnv("APP_NAME", "Share a Ride"),
		Version: getEnv("APP_VERSION", "1.0.0"),
	}

	// Email configuration
	cfg.Email = EmailConfig{
		Host:     getEnv("EMAIL_HOST", "smtp.example.com"),
		Port:     getIntEnv("EMAIL_PORT", 587),
		Username: getEnv("EMAIL_USERNAME", "user@example.com"),
		Password: getEnv("EMAIL_PASSWORD", "password"),
		From:     getEnv("EMAIL_FROM", "noreply@example.com"),
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
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// Returns the full server address
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if str, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(str); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if str, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(str); err == nil {
			return value
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if str, exists := os.LookupEnv(key); exists {
		if value, err := strconv.ParseBool(str); err == nil {
			return value
		}
	}
	return defaultValue
}
