package usecase

import (
	"testing"

	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB is a mock implementation of *gorm.DB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	args := m.Called(dest, conds)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	uc := NewUserUseCase(&gorm.DB{})

	user := &entity.User{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		Role:     entity.RoleUser,
	}

	tests := []struct {
		name    string
		setup   func()
		wantErr error
	}{
		{
			name: "Successful user creation",
			setup: func() {
				mockDB.On("Create", user).Return(&gorm.DB{Error: nil})
			},
			wantErr: nil,
		},
		{
			name: "Failed user creation",
			setup: func() {
				mockDB.On("Create", user).Return(&gorm.DB{Error: gorm.ErrInvalidData})
			},
			wantErr: errs.ErrUserCreation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := uc.CreateUser(user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	mockDB := new(MockDB)
	uc := NewUserUseCase(&gorm.DB{})

	user := &entity.User{
		Model: gorm.Model{ID: 1},
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleUser,
	}

	tests := []struct {
		name    string
		id      uint
		setup   func()
		want    *entity.User
		wantErr error
	}{
		{
			name: "User found",
			id:   1,
			setup: func() {
				mockDB.On("First", mock.Anything, uint(1)).Return(&gorm.DB{Error: nil})
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "User not found",
			id:   2,
			setup: func() {
				mockDB.On("First", mock.Anything, uint(2)).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
			},
			want:    nil,
			wantErr: errs.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := uc.GetUserByID(tt.id)
			assert.Equal(t, tt.wantErr, err)
			if tt.want != nil {
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Email, got.Email)
			}
		})
	}
}

func TestGenerateTokens(t *testing.T) {
	uc := NewUserUseCase(&gorm.DB{})

	user := &entity.User{
		Model: gorm.Model{ID: 1},
		Email: "test@example.com",
		Name:  "Test User",
		Role:  entity.RoleUser,
	}

	tokens, err := uc.GenerateTokens(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestFindOrCreateOAuthUser(t *testing.T) {
	mockDB := new(MockDB)
	uc := NewUserUseCase(&gorm.DB{})

	userInfo := map[string]interface{}{
		"email": "oauth@example.com",
		"name":  "OAuth User",
	}

	tests := []struct {
		name     string
		provider string
		setup    func()
		want     *entity.User
		wantErr  bool
	}{
		{
			name:     "Create new OAuth user",
			provider: "google",
			setup: func() {
				mockDB.On("Where", "email = ?", "oauth@example.com").Return(&gorm.DB{Error: gorm.ErrRecordNotFound})
				mockDB.On("Create", mock.Anything).Return(&gorm.DB{Error: nil})
			},
			want: &entity.User{
				Email:           "oauth@example.com",
				Name:            "OAuth User",
				Role:            entity.RoleUser,
				IsEmailVerified: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := uc.FindOrCreateOAuthUser(tt.provider, userInfo)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Email, got.Email)
				assert.Equal(t, tt.want.Name, got.Name)
				assert.Equal(t, tt.want.IsEmailVerified, got.IsEmailVerified)
			}
		})
	}
}
