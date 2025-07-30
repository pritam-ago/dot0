package models

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	PIN           string `gorm:"primaryKey"`
	CreatedAt     time.Time
	ExpiresAt     time.Time
	LastConnected time.Time
	PCConnected   bool
	UserConnected bool
	BaseDirectory *string
}

func CreateSession(db *gorm.DB, pin string) (*Session, error) {
	log.Printf("Creating session for PIN: %s", pin)

	// Check if it already exists
	var existing Session
	err := db.First(&existing, "pin = ?", pin).Error
	if err == nil {
		log.Printf("PIN %s already exists in database", pin)
		return nil, errors.New("PIN already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Database error checking existing PIN %s: %v", pin, err)
		return nil, err
	}

	log.Printf("PIN %s does not exist, proceeding with creation", pin)

	session := &Session{
		PIN:           pin,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(15 * 24 * time.Hour),
		LastConnected: time.Now(),
		PCConnected:   false,
	}

	log.Printf("Session object created for PIN %s: %+v", pin, session)

	err = db.Create(session).Error
	if err != nil {
		log.Printf("Failed to create session in database for PIN %s: %v", pin, err)
		return nil, err
	}

	log.Printf("Successfully created session in database for PIN: %s", pin)
	return session, nil
}

func GetSession(db *gorm.DB, pin string) (*Session, error) {
	log.Printf("Retrieving session for PIN: %s", pin)

	var s Session
	err := db.First(&s, "pin = ?", pin).Error
	if err != nil {
		log.Printf("Failed to retrieve session for PIN %s: %v", pin, err)
		return &s, err
	}

	log.Printf("Successfully retrieved session for PIN: %s", pin)
	log.Printf("Retrieved session details - PIN: %s, BaseDirectory: %v, UserConnected: %v, PCConnected: %v",
		s.PIN, s.BaseDirectory, s.UserConnected, s.PCConnected)

	return &s, err
}

func UpdatePCConnectionStatus(db *gorm.DB, pin string, connected bool) error {
	return db.Model(&Session{}).
		Where("pin = ?", pin).
		Updates(map[string]interface{}{
			"pc_connected":   connected,
			"last_connected": time.Now(),
		}).Error
}

func UpdateBaseDirectory(db *gorm.DB, pin string, path string) error {
	return db.Model(&Session{}).
		Where("pin = ?", pin).
		Update("base_directory", path).Error
}

func MarkUserConnected(db *gorm.DB, pin string, connected bool) error {
	return db.Model(&Session{}).
		Where("pin = ?", pin).
		Update("user_connected", connected).Error
}
