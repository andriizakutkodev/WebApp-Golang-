package storage

import (
	"webapp/internal/config"
	"webapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	Db *gorm.DB
}

func InitStorage(cfg *config.Config) *Storage {
	var storage Storage

	db, err := gorm.Open(postgres.Open(cfg.ConnectionString), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	storage.Db = db

	return &storage
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(models.User{})
}
