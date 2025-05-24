package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a new Gin engine
	r := gin.New()
	r.GET("/health", HealthCheck)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)

	// Serve the request
	r.ServeHTTP(w, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the response body
	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the response fields
	assert.Equal(t, "ok", response.Status)
	assert.NotEmpty(t, response.Timestamp)
	assert.Equal(t, "1.0.0", response.Version)
	assert.Contains(t, []string{"up", "down"}, response.DBStatus)
}
