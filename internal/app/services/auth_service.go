package services

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/errors"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"github.com/sayeed1999/share-a-ride/internal/domain/repositories"
	"github.com/sayeed1999/share-a-ride/internal/domain/services"
	"github.com/sayeed1999/share-a-ride/internal/provider/token"
)

type authService struct {
	userRepo      repositories.UserRepository
	tokenProvider token.Provider
}

func NewAuthService(userRepo repositories.UserRepository, tokenProvider token.Provider) services.AuthService {
	return &authService{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}
}

func (s *authService) Register(ctx context.Context, input services.RegisterUserInput) (*models.User, *services.TokenPair, error) {
	// Check if email exists
	if _, err := s.userRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, nil, errors.ErrEmailExists
	}

	// Check if phone exists
	if _, err := s.userRepo.FindByPhone(ctx, input.Phone); err == nil {
		return nil, nil, errors.ErrPhoneExists
	}

	// Create new user
	user, err := models.NewUser(input.Name, input.Email, input.Phone, input.Password, input.UserType)
	if err != nil {
		return nil, nil, err
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, &services.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, input services.LoginInput) (*models.User, *services.TokenPair, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, nil, errors.ErrInvalidCredentials
	}

	// Validate password
	if !user.ValidatePassword(input.Password) {
		return nil, nil, errors.ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, err := s.tokenProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, nil, err
	}

	return user, &services.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*services.TokenPair, error) {
	// Validate refresh token
	claims, err := s.tokenProvider.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Generate new tokens
	newAccessToken, err := s.tokenProvider.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.tokenProvider.GenerateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &services.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	// Validate token
	claims, err := s.tokenProvider.ValidateToken(token)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Get user
	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}
