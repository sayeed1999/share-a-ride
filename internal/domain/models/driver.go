package models

import (
	"time"

	"github.com/google/uuid"
)

type VehicleType string

const (
	VehicleTypeCar  VehicleType = "car"
	VehicleTypeBike VehicleType = "bike"
)

type DocumentType string

const (
	DocumentTypeLicense      DocumentType = "license"
	DocumentTypeRegistration DocumentType = "registration"
	DocumentTypeInsurance    DocumentType = "insurance"
)

type Vehicle struct {
	Type        VehicleType `json:"type" gorm:"size:20;not null"`
	Model       string      `json:"model" gorm:"size:100;not null"`
	PlateNumber string      `json:"plate_number" gorm:"size:20;not null"`
}

type Document struct {
	ID        string       `json:"id" gorm:"primaryKey;type:uuid"`
	DriverID  string       `json:"driver_id" gorm:"type:uuid;not null"`
	Type      DocumentType `json:"type" gorm:"size:20;not null"`
	FileURL   string       `json:"file_url" gorm:"size:255;not null"`
	CreatedAt time.Time    `json:"created_at" gorm:"not null"`
}

type Location struct {
	Latitude  float64 `json:"latitude" gorm:"type:decimal(10,8)"`
	Longitude float64 `json:"longitude" gorm:"type:decimal(11,8)"`
}

type Driver struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid"`
	UserID          string     `json:"user_id" gorm:"type:uuid;not null"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	LicenseNumber   string     `json:"license_number" gorm:"size:50;not null;unique"`
	Vehicle         Vehicle    `json:"vehicle" gorm:"embedded"`
	IsVerified      bool       `json:"is_verified" gorm:"default:false"`
	IsAvailable     bool       `json:"is_available" gorm:"default:false"`
	CurrentLocation Location   `json:"current_location" gorm:"embedded"`
	Documents       []Document `json:"documents,omitempty" gorm:"foreignKey:DriverID"`
	CreatedAt       time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"not null"`
}

func NewDriver(userID, licenseNumber string, vehicle Vehicle) *Driver {
	return &Driver{
		ID:            uuid.New().String(),
		UserID:        userID,
		LicenseNumber: licenseNumber,
		Vehicle:       vehicle,
		IsVerified:    false,
		IsAvailable:   false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func NewDocument(driverID string, docType DocumentType, fileURL string) *Document {
	return &Document{
		ID:        uuid.New().String(),
		DriverID:  driverID,
		Type:      docType,
		FileURL:   fileURL,
		CreatedAt: time.Now(),
	}
}

func (d *Driver) UpdateLocation(lat, lng float64) {
	d.CurrentLocation = Location{
		Latitude:  lat,
		Longitude: lng,
	}
	d.UpdatedAt = time.Now()
}

func (d *Driver) UpdateAvailability(isAvailable bool) {
	d.IsAvailable = isAvailable
	d.UpdatedAt = time.Now()
}

func (d *Driver) SetVerified(verified bool) {
	d.IsVerified = verified
	d.UpdatedAt = time.Now()
}
