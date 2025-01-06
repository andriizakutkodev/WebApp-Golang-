package main

import (
	"webapp/internal/config"
	"webapp/internal/logger"
)

func main() {
	// Initialize config from local.yaml file
	cfg := config.GetConfig()
	// Init logger
	log := logger.InitLogger(cfg.Env)
	log.Info("application started")
}
