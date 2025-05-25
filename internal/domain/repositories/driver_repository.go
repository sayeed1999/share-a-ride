package repositories

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/models"
)

type DriverRepository interface {
	Create(ctx context.Context, driver *models.Driver) error
	FindByID(ctx context.Context, id string) (*models.Driver, error)
	FindByUserID(ctx context.Context, userID string) (*models.Driver, error)
	FindByLicenseNumber(ctx context.Context, licenseNumber string) (*models.Driver, error)
	Update(ctx context.Context, driver *models.Driver) error
	Delete(ctx context.Context, id string) error

	// Document related operations
	AddDocument(ctx context.Context, document *models.Document) error
	GetDocuments(ctx context.Context, driverID string) ([]models.Document, error)
	DeleteDocument(ctx context.Context, documentID string) error

	// Location and availability
	UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error
	UpdateAvailability(ctx context.Context, driverID string, isAvailable bool) error
	FindAvailableNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]models.Driver, error)
}
