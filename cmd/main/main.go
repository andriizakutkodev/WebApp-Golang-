package main

import (
	"webapp/internal/config"
	"webapp/internal/logger"
	s "webapp/internal/storage"
)

func main() {
	// Initialize config from local.yaml file
	cfg := config.GetConfig()
	// Init logger
	log := logger.InitLogger(cfg.Env)
	// Init storage and migrate domain models
	storage := s.InitStorage(cfg)
	s.AutoMigrate(storage.Db)

	log.Info("application started")
}
