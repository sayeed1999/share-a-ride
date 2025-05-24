package usecase

import "github.com/sayeed1999/share-a-ride/internal/domain/entity"

// UserUseCase defines the interface for user-related business logic
type UserUseCase interface {
	CreateUser(user *entity.User) error
	GetUserByID(id uint) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByVerifyToken(token string) (*entity.User, error)
	GetUserByResetToken(token string) (*entity.User, error)
	UpdateUser(user *entity.User) error
	FindOrCreateOAuthUser(provider string, userInfo map[string]interface{}) (*entity.User, error)
	GenerateTokens(user *entity.User) (*TokenPair, error)
}
