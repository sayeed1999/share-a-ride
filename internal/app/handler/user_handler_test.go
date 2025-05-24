package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sayeed1999/share-a-ride/internal/app/dto"
	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/errs"
	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/provider/email"
)

// Ensure MockEmailService implements EmailServiceInterface
var _ email.EmailServiceInterface = (*MockEmailService)(nil)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendVerificationEmail(email, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordResetEmail(email, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}

func setupTestRouter(userUseCase usecase.UserUseCase, emailService email.EmailServiceInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewUserHandler(userUseCase, emailService)

	r.POST("/users", handler.CreateUser)
	r.GET("/users/:id", handler.GetUser)
	r.POST("/login", handler.Login)
	r.GET("/verify-email", handler.VerifyEmail)
	r.POST("/reset-password", handler.RequestPasswordReset)
	r.POST("/reset-password/confirm", handler.ResetPassword)

	return r
}

func TestCreateUser(t *testing.T) {
	mockUserUseCase := new(MockUserUseCase)
	mockEmailService := new(MockEmailService)
	router := setupTestRouter(mockUserUseCase, mockEmailService)

	tests := []struct {
		name           string
		request        dto.UserCreateRequest
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful user creation",
			request: dto.UserCreateRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserUseCase.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)
				mockEmailService.On("SendVerificationEmail", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid email format",
			request: dto.UserCreateRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			setupMocks:     func() {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			reqBody, _ := json.Marshal(tt.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUserUseCase.AssertExpectations(t)
			mockEmailService.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	mockUserUseCase := new(MockUserUseCase)
	mockEmailService := new(MockEmailService)
	router := setupTestRouter(mockUserUseCase, mockEmailService)

	validUser := &entity.User{
		Email:    "test@example.com",
		Password: "$2a$10$abcdefghijklmnopqrstuvwxyz", // hashed password
		Name:     "Test User",
	}

	tests := []struct {
		name           string
		request        dto.LoginRequest
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful login",
			request: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func() {
				mockUserUseCase.On("GetUserByEmail", "test@example.com").Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid credentials",
			request: dto.LoginRequest{
				Email:    "wrong@example.com",
				Password: "wrongpass",
			},
			setupMocks: func() {
				mockUserUseCase.On("GetUserByEmail", "wrong@example.com").Return(nil, errs.ErrUserNotFound)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			reqBody, _ := json.Marshal(tt.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUserUseCase.AssertExpectations(t)
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	mockUserUseCase := new(MockUserUseCase)
	mockEmailService := new(MockEmailService)
	router := setupTestRouter(mockUserUseCase, mockEmailService)

	validUser := &entity.User{
		Email:       "test@example.com",
		VerifyToken: "valid-token",
	}

	tests := []struct {
		name           string
		token          string
		setupMocks     func()
		expectedStatus int
		expectedError  string
	}{
		{
			name:  "Successful verification",
			token: "valid-token",
			setupMocks: func() {
				mockUserUseCase.On("GetUserByVerifyToken", "valid-token").Return(validUser, nil)
				mockUserUseCase.On("UpdateUser", mock.AnythingOfType("*entity.User")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Invalid token",
			token: "invalid-token",
			setupMocks: func() {
				mockUserUseCase.On("GetUserByVerifyToken", "invalid-token").Return(nil, errs.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "invalid verification token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/verify-email?token="+tt.token, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.expectedError)
			}

			mockUserUseCase.AssertExpectations(t)
		})
	}
}
