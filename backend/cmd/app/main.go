package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/pritam-ago/dot0/backend/config"
	"github.com/pritam-ago/dot0/backend/db"
	"github.com/pritam-ago/dot0/backend/internal/repositories"
	"github.com/pritam-ago/dot0/backend/internal/routes"
	"github.com/pritam-ago/dot0/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)

	// Initialize router
	router := gin.Default()

	// Setup routes
	routes.SetupAuthRoutes(router, authService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
