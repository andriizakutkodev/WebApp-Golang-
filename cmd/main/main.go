package main

import (
	"webapp/internal/config"
	"webapp/internal/routes"
	"webapp/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config from local.yaml file
	cfg := config.GetConfig()
	// Init storage and migrate domain models
	db := storage.InitStorage(cfg)
	// Init router and map all routes
	r := gin.Default()

	routes.RegisterRoutes(r, db, cfg)

	r.Run(":8080")
}
