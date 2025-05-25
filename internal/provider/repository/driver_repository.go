package repository

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/errors"
	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"github.com/sayeed1999/share-a-ride/internal/domain/repositories"
	"gorm.io/gorm"
)

type driverRepository struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) repositories.DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) Create(ctx context.Context, driver *models.Driver) error {
	return r.db.WithContext(ctx).Create(driver).Error
}

func (r *driverRepository) FindByID(ctx context.Context, id string) (*models.Driver, error) {
	var driver models.Driver
	if err := r.db.WithContext(ctx).Preload("Documents").First(&driver, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrDriverNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) FindByUserID(ctx context.Context, userID string) (*models.Driver, error) {
	var driver models.Driver
	if err := r.db.WithContext(ctx).Preload("Documents").First(&driver, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrDriverNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) FindByLicenseNumber(ctx context.Context, licenseNumber string) (*models.Driver, error) {
	var driver models.Driver
	if err := r.db.WithContext(ctx).First(&driver, "license_number = ?", licenseNumber).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrDriverNotFound
		}
		return nil, err
	}
	return &driver, nil
}

func (r *driverRepository) Update(ctx context.Context, driver *models.Driver) error {
	return r.db.WithContext(ctx).Save(driver).Error
}

func (r *driverRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Driver{}, "id = ?", id).Error
}

func (r *driverRepository) AddDocument(ctx context.Context, document *models.Document) error {
	return r.db.WithContext(ctx).Create(document).Error
}

func (r *driverRepository) GetDocuments(ctx context.Context, driverID string) ([]models.Document, error) {
	var documents []models.Document
	if err := r.db.WithContext(ctx).Where("driver_id = ?", driverID).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

func (r *driverRepository) DeleteDocument(ctx context.Context, documentID string) error {
	return r.db.WithContext(ctx).Delete(&models.Document{}, "id = ?", documentID).Error
}

func (r *driverRepository) UpdateLocation(ctx context.Context, driverID string, lat, lng float64) error {
	return r.db.WithContext(ctx).Model(&models.Driver{}).
		Where("id = ?", driverID).
		Updates(map[string]interface{}{
			"current_latitude":  lat,
			"current_longitude": lng,
			"updated_at":        gorm.Expr("NOW()"),
		}).Error
}

func (r *driverRepository) UpdateAvailability(ctx context.Context, driverID string, isAvailable bool) error {
	return r.db.WithContext(ctx).Model(&models.Driver{}).
		Where("id = ?", driverID).
		Updates(map[string]interface{}{
			"is_available": isAvailable,
			"updated_at":   gorm.Expr("NOW()"),
		}).Error
}

func (r *driverRepository) FindAvailableNearby(ctx context.Context, lat, lng float64, radiusKm float64) ([]models.Driver, error) {
	var drivers []models.Driver

	// Using the Haversine formula to calculate distance
	query := `
		SELECT *,
		(6371 * acos(cos(radians(?)) * 
		cos(radians(current_latitude)) * 
		cos(radians(current_longitude) - 
		radians(?)) + 
		sin(radians(?)) * 
		sin(radians(current_latitude)))) AS distance 
		FROM drivers 
		WHERE is_available = true 
		AND is_verified = true
		HAVING distance <= ? 
		ORDER BY distance
	`

	if err := r.db.WithContext(ctx).Raw(query, lat, lng, lat, radiusKm).Scan(&drivers).Error; err != nil {
		return nil, err
	}

	return drivers, nil
}
