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
	s := storage.InitStorage(cfg)
	storage.AutoMigrate(s.Db)
	// Init router and map all routes
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
