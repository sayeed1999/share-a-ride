package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/domain/errors"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"github.com/sayeed1999/share-a-ride/internal/domain/services"
)

type verifyDriverRequest struct {
	LicenseNumber string            `json:"license_number" binding:"required"`
	Vehicle       vehicleRequest    `json:"vehicle" binding:"required"`
	Documents     []documentRequest `json:"documents" binding:"required,min=1"`
}

type vehicleRequest struct {
	Type        models.VehicleType `json:"type" binding:"required,oneof=car bike"`
	Model       string             `json:"model" binding:"required"`
	PlateNumber string             `json:"plate_number" binding:"required"`
}

type documentRequest struct {
	Type    models.DocumentType `json:"type" binding:"required,oneof=license registration insurance"`
	FileURL string              `json:"file_url" binding:"required,url"`
}

type updateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type updateAvailabilityRequest struct {
	IsAvailable bool `json:"is_available" binding:"required"`
}

type DriverHandler struct {
	driverService services.DriverService
}

func NewDriverHandler(driverService services.DriverService) *DriverHandler {
	return &DriverHandler{
		driverService: driverService,
	}
}

func (h *DriverHandler) VerifyDriver(c *gin.Context) {
	var req verifyDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)

	// Convert request to service input
	documents := make([]services.DocumentInput, len(req.Documents))
	for i, doc := range req.Documents {
		documents[i] = services.DocumentInput{
			Type:    doc.Type,
			FileURL: doc.FileURL,
		}
	}

	driver, err := h.driverService.VerifyDriver(c.Request.Context(), user.ID, services.VerifyDriverInput{
		LicenseNumber: req.LicenseNumber,
		Vehicle: models.Vehicle{
			Type:        req.Vehicle.Type,
			Model:       req.Vehicle.Model,
			PlateNumber: req.Vehicle.PlateNumber,
		},
		Documents: documents,
	})

	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case errors.ErrDriverExists, errors.ErrLicenseExists:
			status = http.StatusConflict
		case errors.ErrUnauthorizedAccess:
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"data": gin.H{
			"verification_id": driver.ID,
			"status":          "pending",
			"submitted_at":    driver.CreatedAt,
		},
	})
}

func (h *DriverHandler) UpdateLocation(c *gin.Context) {
	var req updateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	driver, err := h.driverService.GetDriverByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	err = h.driverService.UpdateLocation(c.Request.Context(), driver.ID, services.UpdateLocationInput{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"driver_id":        driver.ID,
			"current_location": req,
			"updated_at":       driver.UpdatedAt,
		},
	})
}

func (h *DriverHandler) UpdateAvailability(c *gin.Context) {
	var req updateAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(*models.User)
	driver, err := h.driverService.GetDriverByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	err = h.driverService.UpdateAvailability(c.Request.Context(), driver.ID, req.IsAvailable)
	if err != nil {
		status := http.StatusInternalServerError
		if err == errors.ErrDriverNotVerified {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"driver_id":    driver.ID,
			"is_available": req.IsAvailable,
			"updated_at":   driver.UpdatedAt,
		},
	})
}

func (h *DriverHandler) GetProfile(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	driver, err := h.driverService.GetDriverByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    driver,
	})
}

func (h *DriverHandler) GetDocuments(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	driver, err := h.driverService.GetDriverByUserID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	documents, err := h.driverService.GetDocuments(c.Request.Context(), driver.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    documents,
	})
}
