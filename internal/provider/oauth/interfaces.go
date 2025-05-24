package oauth

import "golang.org/x/oauth2"

// OAuthService defines the interface for OAuth operations
type OAuthService interface {
	RegisterProvider(name, clientID, clientSecret, redirectURL string, scopes []string)
	GetAuthURL(provider, state string) (string, error)
	Exchange(provider, code string) (*oauth2.Token, error)
	GetUserInfo(provider string, token *oauth2.Token) (map[string]interface{}, error)
}
