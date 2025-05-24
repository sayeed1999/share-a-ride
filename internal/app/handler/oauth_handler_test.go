package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"

	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/mocks"
	"github.com/sayeed1999/share-a-ride/internal/provider/oauth"
)

func setupOAuthTestRouter(userUseCase *mocks.MockUserUseCase, oauthService *mocks.MockOAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	handler := NewOAuthHandler(userUseCase, oauthService)

	r.GET("/oauth/:provider", handler.InitiateOAuth)
	r.GET("/oauth/:provider/callback", handler.OAuthCallback)

	return r
}

func TestInitiateOAuth(t *testing.T) {
	mockUserUseCase := new(mocks.MockUserUseCase)
	mockOAuthService := new(mocks.MockOAuthService)
	router := setupOAuthTestRouter(mockUserUseCase, mockOAuthService)

	tests := []struct {
		name           string
		provider       string
		state          string
		setupMocks     func()
		expectedStatus int
		expectedURL    string
	}{
		{
			name:     "Successful OAuth initiation",
			provider: "google",
			state:    "random-state",
			setupMocks: func() {
				mockOAuthService.On("GetAuthURL", "google", "random-state").Return(
					"https://accounts.google.com/o/oauth2/auth?client_id=123&state=random-state",
					nil,
				)
			},
			expectedStatus: http.StatusTemporaryRedirect,
			expectedURL:    "https://accounts.google.com/o/oauth2/auth?client_id=123&state=random-state",
		},
		{
			name:     "Invalid provider",
			provider: "invalid",
			state:    "random-state",
			setupMocks: func() {
				mockOAuthService.On("GetAuthURL", "invalid", "random-state").Return(
					"",
					oauth.ErrInvalidProvider,
				)
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/oauth/"+tt.provider+"?state="+tt.state, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedURL != "" {
				assert.Equal(t, tt.expectedURL, w.Header().Get("Location"))
			}

			mockOAuthService.AssertExpectations(t)
		})
	}
}

func TestOAuthCallback(t *testing.T) {
	mockUserUseCase := new(mocks.MockUserUseCase)
	mockOAuthService := new(mocks.MockOAuthService)
	router := setupOAuthTestRouter(mockUserUseCase, mockOAuthService)

	validToken := &oauth2.Token{
		AccessToken: "valid-access-token",
	}

	validUserInfo := map[string]interface{}{
		"email": "test@example.com",
		"name":  "Test User",
	}

	validUser := &entity.User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	tests := []struct {
		name           string
		provider       string
		code           string
		state          string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name:     "Successful OAuth callback",
			provider: "google",
			code:     "valid-code",
			state:    "valid-state",
			setupMocks: func() {
				mockOAuthService.On("Exchange", "google", "valid-code").Return(validToken, nil)
				mockOAuthService.On("GetUserInfo", "google", validToken).Return(validUserInfo, nil)
				mockUserUseCase.On("FindOrCreateOAuthUser", "google", validUserInfo).Return(validUser, nil)
				mockUserUseCase.On("GenerateTokens", validUser).Return(&usecase.TokenPair{
					AccessToken:  "jwt-token",
					RefreshToken: "refresh-token",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Invalid code",
			provider: "google",
			code:     "invalid-code",
			state:    "valid-state",
			setupMocks: func() {
				mockOAuthService.On("Exchange", "google", "invalid-code").Return(nil, oauth.ErrInvalidCode)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/oauth/"+tt.provider+"/callback?code="+tt.code+"&state="+tt.state, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockOAuthService.AssertExpectations(t)
			mockUserUseCase.AssertExpectations(t)
		})
	}
}
