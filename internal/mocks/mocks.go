package mocks

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/provider/oauth"
)

// MockUserUseCase is a mock implementation of usecase.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

// Ensure MockUserUseCase implements usecase.UserUseCase
var _ usecase.UserUseCase = (*MockUserUseCase)(nil)

func (m *MockUserUseCase) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) GetUserByID(id uint) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByVerifyToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GetUserByResetToken(token string) (*entity.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserUseCase) FindOrCreateOAuthUser(provider string, userInfo map[string]interface{}) (*entity.User, error) {
	args := m.Called(provider, userInfo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUseCase) GenerateTokens(user *entity.User) (*usecase.TokenPair, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TokenPair), args.Error(1)
}

// MockOAuthService is a mock implementation of oauth.OAuthService
type MockOAuthService struct {
	mock.Mock
}

// Ensure MockOAuthService implements oauth.OAuthService
var _ oauth.OAuthService = (*MockOAuthService)(nil)

func (m *MockOAuthService) GetAuthURL(provider, state string) (string, error) {
	args := m.Called(provider, state)
	return args.String(0), args.Error(1)
}

func (m *MockOAuthService) Exchange(provider, code string) (*oauth2.Token, error) {
	args := m.Called(provider, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockOAuthService) GetUserInfo(provider string, token *oauth2.Token) (map[string]interface{}, error) {
	args := m.Called(provider, token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockOAuthService) RegisterProvider(name, clientID, clientSecret, redirectURL string, scopes []string) {
	m.Called(name, clientID, clientSecret, redirectURL, scopes)
}
