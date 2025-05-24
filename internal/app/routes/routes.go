package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/app/handler"
	"github.com/sayeed1999/share-a-ride/internal/app/middleware"
	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/sayeed1999/share-a-ride/internal/domain/usecase"
	"github.com/sayeed1999/share-a-ride/internal/provider/db"
	"github.com/sayeed1999/share-a-ride/internal/provider/email"
	"github.com/sayeed1999/share-a-ride/internal/provider/oauth"
)

type RouterConfig struct {
	EnableOAuth    bool
	OAuthProviders map[string]OAuthProviderConfig
}

type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

func SetupRoutes(router *gin.Engine, cfg *config.Config, routerCfg *RouterConfig) {
	// Initialize dependencies
	userUseCase := usecase.NewUserUseCase(db.DB)
	emailService := email.NewEmailService(cfg)
	userHandler := handler.NewUserHandler(userUseCase, emailService)

	// Initialize OAuth service if enabled
	var oauthHandler *handler.OAuthHandler
	if routerCfg != nil && routerCfg.EnableOAuth {
		oauthService := oauth.NewOAuthService()

		// Register configured providers
		for provider, providerCfg := range routerCfg.OAuthProviders {
			oauthService.RegisterProvider(
				provider,
				providerCfg.ClientID,
				providerCfg.ClientSecret,
				providerCfg.RedirectURL,
				providerCfg.Scopes,
			)
		}

		oauthHandler = handler.NewOAuthHandler(userUseCase, oauthService)
	}

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimiter())

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.CreateUser)
			auth.POST("/login", userHandler.Login)
			auth.GET("/verify-email", userHandler.VerifyEmail)
			auth.POST("/forgot-password", userHandler.RequestPasswordReset)
			auth.POST("/reset-password", userHandler.ResetPassword)

			// OAuth routes (if enabled)
			if oauthHandler != nil {
				oauth := auth.Group("/oauth")
				{
					oauth.GET("/:provider", oauthHandler.InitiateOAuth)
					oauth.GET("/:provider/callback", oauthHandler.OAuthCallback)
				}
			}
		}

		// Protected routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/:id", middleware.IsUser(), userHandler.GetUser)
		}
	}
}
