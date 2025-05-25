package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/app/http/handlers"
	"github.com/sayeed1999/share-a-ride/internal/app/http/middleware"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handlers.AuthHandler
	driverHandler  *handlers.DriverHandler
	authMiddleware *middleware.AuthMiddleware
}

func New(authHandler *handlers.AuthHandler, driverHandler *handlers.DriverHandler, authMiddleware *middleware.AuthMiddleware) *Router {
	r := &Router{
		engine:         gin.Default(),
		authHandler:    authHandler,
		driverHandler:  driverHandler,
		authMiddleware: authMiddleware,
	}
	return r
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}

func (r *Router) SetupRoutes() {
	// Auth routes
	auth := r.engine.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}

	// Driver routes
	drivers := r.engine.Group("/drivers")
	drivers.Use(r.authMiddleware.Authenticate())
	{
		drivers.POST("/verify", r.authMiddleware.RequireDriver(), r.driverHandler.VerifyDriver)
		drivers.PUT("/location", r.authMiddleware.RequireDriver(), r.driverHandler.UpdateLocation)
		drivers.PUT("/availability", r.authMiddleware.RequireDriver(), r.driverHandler.UpdateAvailability)
		drivers.GET("/profile", r.authMiddleware.RequireDriver(), r.driverHandler.GetProfile)
		drivers.GET("/documents", r.authMiddleware.RequireDriver(), r.driverHandler.GetDocuments)
	}
}
