package storage

import (
	"webapp/internal/config"
	"webapp/internal/models"

	"gorm.io/driver/postgres"
	g "gorm.io/gorm"
)

type Storage struct {
	Db *g.DB
}

func InitStorage(cfg *config.Config) *Storage {
	var storage Storage

	db, err := g.Open(postgres.Open(cfg.ConnectionString), &g.Config{})

	if err != nil {
		panic(err.Error())
	}

	storage.Db = db

	return &storage
}

func AutoMigrate(db *g.DB) {
	db.AutoMigrate(models.User{})
}
