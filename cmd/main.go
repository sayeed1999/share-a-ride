package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r.Engine(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
