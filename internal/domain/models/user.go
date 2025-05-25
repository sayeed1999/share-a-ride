package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	UserTypeRider  UserType = "rider"
	UserTypeDriver UserType = "driver"
)

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:255;not null;unique"`
	Phone     string    `json:"phone" gorm:"size:20;not null;unique"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	UserType  UserType  `json:"user_type" gorm:"size:10;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func NewUser(name, email, phone, password string, userType UserType) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Password:  string(hashedPassword),
		UserType:  userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsDriver() bool {
	return u.UserType == UserTypeDriver
}

func (u *User) IsRider() bool {
	return u.UserType == UserTypeRider
}
