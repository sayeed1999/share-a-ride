package oauth

import (
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Provider represents an OAuth2 provider
type Provider struct {
	Name         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Config       *oauth2.Config
}

// OAuthServiceImpl implements the OAuthService interface
type OAuthServiceImpl struct {
	providers map[string]*Provider
}

// NewOAuthService creates a new OAuth service
func NewOAuthService() OAuthService {
	return &OAuthServiceImpl{
		providers: make(map[string]*Provider),
	}
}

// RegisterProvider registers a new OAuth2 provider
func (s *OAuthServiceImpl) RegisterProvider(name, clientID, clientSecret, redirectURL string, scopes []string) {
	var config *oauth2.Config

	switch name {
	case "google":
		config = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       scopes,
			Endpoint:     google.Endpoint,
		}
		// Add more providers as needed
		// case "github":
		//     config = &oauth2.Config{...}
	}

	s.providers[name] = &Provider{
		Name:         name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Config:       config,
	}
}

// GetAuthURL returns the OAuth2 authorization URL
func (s *OAuthServiceImpl) GetAuthURL(provider string, state string) (string, error) {
	p, exists := s.providers[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not found", provider)
	}

	return p.Config.AuthCodeURL(state), nil
}

// Exchange exchanges the authorization code for tokens
func (s *OAuthServiceImpl) Exchange(provider, code string) (*oauth2.Token, error) {
	p, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", provider)
	}

	return p.Config.Exchange(oauth2.NoContext, code)
}

// GetUserInfo fetches user information from the OAuth2 provider
func (s *OAuthServiceImpl) GetUserInfo(provider string, token *oauth2.Token) (map[string]interface{}, error) {
	p, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", provider)
	}

	client := p.Config.Client(oauth2.NoContext, token)

	var userInfoURL string
	switch provider {
	case "google":
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	// Add more providers as needed
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return userInfo, nil
}
