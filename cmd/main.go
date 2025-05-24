package main

import (
	"log"

	"github.com/sayeed1999/share-a-ride/internal/app/routes"
	"github.com/sayeed1999/share-a-ride/internal/config"
	"github.com/sayeed1999/share-a-ride/internal/provider/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Initialize database
	if err := db.InitDB(cfg); err != nil {
		log.Fatal("Error initializing database:", err)
	}

	// Create Gin router
	router := gin.Default()

	// Setup routes with OAuth config
	routerCfg := &routes.RouterConfig{
		EnableOAuth:    false, // Set to true to enable OAuth
		OAuthProviders: map[string]routes.OAuthProviderConfig{},
	}
	routes.SetupRoutes(router, cfg, routerCfg)

	// Start server
	log.Printf("Server starting on %s", cfg.Server.GetAddress())
	if err := router.Run(cfg.Server.GetAddress()); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
