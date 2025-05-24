package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestRegisterProvider(t *testing.T) {
	service := NewOAuthService().(*OAuthServiceImpl)

	tests := []struct {
		name         string
		provider     string
		clientID     string
		clientSecret string
		redirectURL  string
		scopes       []string
	}{
		{
			name:         "Register Google provider",
			provider:     "google",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			redirectURL:  "http://localhost:8080/callback",
			scopes:       []string{"email", "profile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.RegisterProvider(tt.provider, tt.clientID, tt.clientSecret, tt.redirectURL, tt.scopes)

			provider, exists := service.providers[tt.provider]
			assert.True(t, exists)
			assert.Equal(t, tt.provider, provider.Name)
			assert.Equal(t, tt.clientID, provider.ClientID)
			assert.Equal(t, tt.clientSecret, provider.ClientSecret)
			assert.Equal(t, tt.redirectURL, provider.RedirectURL)
			assert.Equal(t, tt.scopes, provider.Scopes)
			assert.NotNil(t, provider.Config)
		})
	}
}

func TestGetAuthURL(t *testing.T) {
	service := NewOAuthService()
	service.RegisterProvider(
		"google",
		"test-client-id",
		"test-client-secret",
		"http://localhost:8080/callback",
		[]string{"email", "profile"},
	)

	tests := []struct {
		name        string
		provider    string
		state       string
		wantErr     bool
		wantURLPart string
	}{
		{
			name:        "Valid provider",
			provider:    "google",
			state:       "test-state",
			wantErr:     false,
			wantURLPart: "https://accounts.google.com/o/oauth2/auth",
		},
		{
			name:     "Invalid provider",
			provider: "invalid",
			state:    "test-state",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := service.GetAuthURL(tt.provider, tt.state)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, url, tt.wantURLPart)
			assert.Contains(t, url, tt.state)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	service := NewOAuthService()
	service.RegisterProvider(
		"google",
		"test-client-id",
		"test-client-secret",
		"http://localhost:8080/callback",
		[]string{"email", "profile"},
	)

	// Create a test server to mock Google's user info endpoint
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		}
		json.NewEncoder(w).Encode(userInfo)
	}))
	defer testServer.Close()

	tests := []struct {
		name     string
		provider string
		token    *oauth2.Token
		wantErr  bool
	}{
		{
			name:     "Valid provider and token",
			provider: "google",
			token: &oauth2.Token{
				AccessToken: "test-access-token",
			},
			wantErr: false,
		},
		{
			name:     "Invalid provider",
			provider: "invalid",
			token: &oauth2.Token{
				AccessToken: "test-access-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userInfo, err := service.GetUserInfo(tt.provider, tt.token)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, userInfo)
			assert.Equal(t, "test@example.com", userInfo["email"])
			assert.Equal(t, "Test User", userInfo["name"])
		})
	}
}
