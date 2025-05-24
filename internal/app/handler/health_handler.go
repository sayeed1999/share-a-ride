package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/provider/db"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
	DBStatus  string `json:"db_status"`
}

// HealthCheck handles the health check endpoint
func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := db.DB.DB()
	dbStatus := "up"
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "down"
	}

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0", // This should come from config
		DBStatus:  dbStatus,
	}

	c.JSON(http.StatusOK, response)
}
