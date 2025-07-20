package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pritam-ago/go-relay/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:dhanus17@localhost:5432/dot0"
	}

	// Configure GORM to suppress PostgreSQL logs
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Suppress all GORM logs
	}

	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), gormConfig)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Session{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL database")
}
