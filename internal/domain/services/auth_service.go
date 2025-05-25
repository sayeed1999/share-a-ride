package services

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/models"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
	UserType models.UserType
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService interface {
	Register(ctx context.Context, input RegisterUserInput) (*models.User, *TokenPair, error)
	Login(ctx context.Context, input LoginInput) (*models.User, *TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
