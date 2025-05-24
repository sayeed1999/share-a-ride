package entity

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

type User struct {
	gorm.Model
	Email           string    `gorm:"uniqueIndex;not null" json:"email"`
	Password        string    `json:"-"`
	Name            string    `json:"name"`
	Role            Role      `gorm:"type:varchar(20);default:'user'" json:"role"`
	IsEmailVerified bool      `gorm:"default:false" json:"is_email_verified"`
	VerifyToken     string    `gorm:"size:100" json:"-"`
	ResetToken      string    `gorm:"size:100" json:"-"`
	ResetExpires    time.Time `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
