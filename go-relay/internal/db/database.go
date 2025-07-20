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
	log.Printf("🔍 Initializing database connection...")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️ Warning: .env file not found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgresql://postgres:dhanus17@localhost:5432/dot0"
		log.Printf("🔍 Using default database URL: %s", databaseURL)
	} else {
		log.Printf("🔍 Using database URL from environment: %s", databaseURL)
	}

	// Configure GORM to show errors but suppress other logs
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error), // Show errors but suppress other logs
	}

	var err error
	log.Printf("🔍 Attempting to connect to database...")
	DB, err = gorm.Open(postgres.Open(databaseURL), gormConfig)
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
		log.Fatal("failed to connect to database:", err)
	}

	log.Printf("✅ Successfully connected to PostgreSQL database")
	log.Printf("🔍 Running database migrations...")

	err = DB.AutoMigrate(&models.Session{})
	if err != nil {
		log.Printf("❌ Failed to migrate database: %v", err)
		log.Fatal("failed to migrate database:", err)
	}

	log.Printf("✅ Database migrations completed successfully")
	log.Printf("✅ Database initialization complete")
}
