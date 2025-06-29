package main

import (
	"context"
	"log"

	"github.com/sayeed1999/share-a-ride/internal/app/http/handlers"
	"github.com/sayeed1999/share-a-ride/internal/app/http/middleware"
	"github.com/sayeed1999/share-a-ride/internal/app/http/router"
	"github.com/sayeed1999/share-a-ride/internal/app/services"
	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/sayeed1999/share-a-ride/internal/provider/database"
	"github.com/sayeed1999/share-a-ride/internal/provider/repository"
	"github.com/sayeed1999/share-a-ride/internal/provider/token"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(database.Config{
		DSN: cfg.Database.GetDSN(),
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(context.Background()); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB())
	driverRepo := repository.NewDriverRepository(db.DB())

	// Initialize token provider
	tokenProvider := token.NewJWTProvider(cfg.JWT)

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenProvider)
	driverService := services.NewDriverService(driverRepo, userRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	driverHandler := handlers.NewDriverHandler(driverService)

	// Setup router
	r := router.New(authHandler, driverHandler, authMiddleware)
	r.SetupRoutes()

	// Start Gin server (simple way)
	if err := r.Engine().Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
