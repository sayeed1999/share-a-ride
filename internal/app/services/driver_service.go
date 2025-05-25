package services

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/errors"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"github.com/sayeed1999/share-a-ride/internal/domain/repositories"
	"github.com/sayeed1999/share-a-ride/internal/domain/services"
)

type driverService struct {
	driverRepo repositories.DriverRepository
	userRepo   repositories.UserRepository
}

func NewDriverService(driverRepo repositories.DriverRepository, userRepo repositories.UserRepository) services.DriverService {
	return &driverService{
		driverRepo: driverRepo,
		userRepo:   userRepo,
	}
}

func (s *driverService) VerifyDriver(ctx context.Context, userID string, input services.VerifyDriverInput) (*models.Driver, error) {
	// Check if user exists and is a driver
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	if !user.IsDriver() {
		return nil, errors.ErrUnauthorizedAccess
	}

	// Check if driver already exists
	if _, err := s.driverRepo.FindByUserID(ctx, userID); err == nil {
		return nil, errors.ErrDriverExists
	}

	// Check if license number is already registered
	if _, err := s.driverRepo.FindByLicenseNumber(ctx, input.LicenseNumber); err == nil {
		return nil, errors.ErrLicenseExists
	}

	// Create driver
	driver := models.NewDriver(userID, input.LicenseNumber, input.Vehicle)

	// Create driver documents
	for _, doc := range input.Documents {
		document := models.NewDocument(driver.ID, doc.Type, doc.FileURL)
		driver.Documents = append(driver.Documents, *document)
	}

	// Save driver
	if err := s.driverRepo.Create(ctx, driver); err != nil {
		return nil, err
	}

	return driver, nil
}

func (s *driverService) UpdateLocation(ctx context.Context, driverID string, input services.UpdateLocationInput) error {
	// Validate driver
	driver, err := s.driverRepo.FindByID(ctx, driverID)
	if err != nil {
		return errors.ErrDriverNotFound
	}

	// Update location
	if err := s.driverRepo.UpdateLocation(ctx, driver.ID, input.Latitude, input.Longitude); err != nil {
		return err
	}

	return nil
}

func (s *driverService) UpdateAvailability(ctx context.Context, driverID string, isAvailable bool) error {
	// Validate driver
	driver, err := s.driverRepo.FindByID(ctx, driverID)
	if err != nil {
		return errors.ErrDriverNotFound
	}

	// Check if driver is verified
	if !driver.IsVerified {
		return errors.ErrDriverNotVerified
	}

	// Update availability
	if err := s.driverRepo.UpdateAvailability(ctx, driver.ID, isAvailable); err != nil {
		return err
	}

	return nil
}

func (s *driverService) GetDriverProfile(ctx context.Context, driverID string) (*models.Driver, error) {
	driver, err := s.driverRepo.FindByID(ctx, driverID)
	if err != nil {
		return nil, errors.ErrDriverNotFound
	}
	return driver, nil
}

func (s *driverService) GetDriverByUserID(ctx context.Context, userID string) (*models.Driver, error) {
	driver, err := s.driverRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.ErrDriverNotFound
	}
	return driver, nil
}

func (s *driverService) AddDocument(ctx context.Context, driverID string, input services.DocumentInput) error {
	// Validate driver
	driver, err := s.driverRepo.FindByID(ctx, driverID)
	if err != nil {
		return errors.ErrDriverNotFound
	}

	// Create and save document
	document := models.NewDocument(driver.ID, input.Type, input.FileURL)
	if err := s.driverRepo.AddDocument(ctx, document); err != nil {
		return err
	}

	return nil
}

func (s *driverService) GetDocuments(ctx context.Context, driverID string) ([]models.Document, error) {
	// Validate driver
	if _, err := s.driverRepo.FindByID(ctx, driverID); err != nil {
		return nil, errors.ErrDriverNotFound
	}

	// Get documents
	documents, err := s.driverRepo.GetDocuments(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (s *driverService) DeleteDocument(ctx context.Context, driverID string, documentID string) error {
	// Validate driver
	if _, err := s.driverRepo.FindByID(ctx, driverID); err != nil {
		return errors.ErrDriverNotFound
	}

	// Delete document
	if err := s.driverRepo.DeleteDocument(ctx, documentID); err != nil {
		return err
	}

	return nil
}
