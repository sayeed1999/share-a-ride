package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/provider/oauth"
)

type OAuthHandler struct {
	userUseCase  usecase.UserUseCase
	oauthService oauth.OAuthService
}

func NewOAuthHandler(userUseCase usecase.UserUseCase, oauthService oauth.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		userUseCase:  userUseCase,
		oauthService: oauthService,
	}
}

// InitiateOAuth starts the OAuth2 flow
func (h *OAuthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	state := c.Query("state") // Optional state parameter

	authURL, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OAuthCallback handles the OAuth2 callback
func (h *OAuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	// Verify state if provided
	if expectedState := c.Query("state"); state != expectedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	// Exchange code for token
	token, err := h.oauthService.Exchange(provider, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	// Get user info from provider
	userInfo, err := h.oauthService.GetUserInfo(provider, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// Find or create user
	user, err := h.userUseCase.FindOrCreateOAuthUser(provider, userInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process user"})
		return
	}

	// Generate JWT token
	tokenPair, err := h.userUseCase.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
