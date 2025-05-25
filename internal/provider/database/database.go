package database

import (
	"context"

	"github.com/sayeed1999/share-a-ride/internal/domain/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Provider interface {
	DB() *gorm.DB
	Close() error
	AutoMigrate(ctx context.Context) error
}

type Config struct {
	DSN string
}

type provider struct {
	db *gorm.DB
}

func New(cfg Config) (Provider, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &provider{db: db}, nil
}

func (p *provider) DB() *gorm.DB {
	return p.db
}

func (p *provider) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (p *provider) AutoMigrate(ctx context.Context) error {
	return p.db.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.Driver{},
		&models.Document{},
	)
}
