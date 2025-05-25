package services

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/models"
)

type VerifyDriverInput struct {
	LicenseNumber string
	Vehicle       models.Vehicle
	Documents     []DocumentInput
}

type DocumentInput struct {
	Type    models.DocumentType
	FileURL string
}

type UpdateLocationInput struct {
	Latitude  float64
	Longitude float64
}

type DriverService interface {
	VerifyDriver(ctx context.Context, userID string, input VerifyDriverInput) (*models.Driver, error)
	UpdateLocation(ctx context.Context, driverID string, input UpdateLocationInput) error
	UpdateAvailability(ctx context.Context, driverID string, isAvailable bool) error
	GetDriverProfile(ctx context.Context, driverID string) (*models.Driver, error)
	GetDriverByUserID(ctx context.Context, userID string) (*models.Driver, error)

	// Document management
	AddDocument(ctx context.Context, driverID string, input DocumentInput) error
	GetDocuments(ctx context.Context, driverID string) ([]models.Document, error)
	DeleteDocument(ctx context.Context, driverID string, documentID string) error
}
