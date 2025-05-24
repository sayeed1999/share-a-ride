package usecase

import (
	"fmt"

	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/errs"
	"github.com/sayeed1999/share-a-ride/internal/pkg/jwtutil"

	"gorm.io/gorm"
)

type UserUseCaseImpl struct {
	db *gorm.DB
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserUseCase(db *gorm.DB) UserUseCase {
	return &UserUseCaseImpl{db: db}
}

func (uc *UserUseCaseImpl) CreateUser(user *entity.User) error {
	result := uc.db.Create(user)
	if result.Error != nil {
		return errs.ErrUserCreation
	}
	return nil
}

func (uc *UserUseCaseImpl) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	result := uc.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (uc *UserUseCaseImpl) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	result := uc.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errs.ErrUserNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindOrCreateOAuthUser finds or creates a user from OAuth provider data
func (uc *UserUseCaseImpl) FindOrCreateOAuthUser(provider string, userInfo map[string]interface{}) (*entity.User, error) {
	email, ok := userInfo["email"].(string)
	if !ok {
		return nil, fmt.Errorf("email not found in OAuth user info")
	}

	var user entity.User
	result := uc.db.Where("email = ?", email).First(&user)

	if result.Error == nil {
		// User exists, return it
		return &user, nil
	}

	if result.Error != gorm.ErrRecordNotFound {
		// Unexpected error
		return nil, result.Error
	}

	// User doesn't exist, create new one
	name, _ := userInfo["name"].(string)
	if name == "" {
		name = email // Use email as name if not provided
	}

	user = entity.User{
		Email:           email,
		Name:            name,
		Role:            entity.RoleUser,
		IsEmailVerified: true, // OAuth users are considered verified
	}

	if err := uc.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GenerateTokens generates a new token pair for a user
func (uc *UserUseCaseImpl) GenerateTokens(user *entity.User) (*TokenPair, error) {
	accessToken, err := jwtutil.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwtutil.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *UserUseCaseImpl) GetUserByVerifyToken(token string) (*entity.User, error) {
	var user entity.User
	result := uc.db.Where("verify_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (uc *UserUseCaseImpl) GetUserByResetToken(token string) (*entity.User, error) {
	var user entity.User
	result := uc.db.Where("reset_token = ?", token).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (uc *UserUseCaseImpl) UpdateUser(user *entity.User) error {
	return uc.db.Save(user).Error
}
