package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"github.com/pritam-ago/go-relay/internal/models"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("sessions.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	err = DB.AutoMigrate(&models.Session{})
	if err != nil {
		log.Fatal("failed to migrate db:", err)
	}
}
