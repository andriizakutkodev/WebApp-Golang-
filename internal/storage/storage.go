package storage

import (
	"webapp/internal/config"
	"webapp/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitStorage(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.ConnectionString), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(models.User{}, models.Note{})

	return db
}
