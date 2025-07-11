package models

import (
	"errors"
	"time"
	"gorm.io/gorm"
)

type Session struct {
	PIN            string    `gorm:"primaryKey"`
	CreatedAt      time.Time
	ExpiresAt      time.Time
	LastConnected  time.Time
	PCConnected    bool
}

func CreateSession(db *gorm.DB, pin string) (*Session, error) {
	// Check if it already exists
	var existing Session
	err := db.First(&existing, "pin = ?", pin).Error
	if err == nil {
		return nil, errors.New("PIN already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	session := &Session{
		PIN:           pin,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(15 * 24 * time.Hour),
		LastConnected: time.Now(),
		PCConnected:   false,
	}
	return session, db.Create(session).Error
}

func GetSession(db *gorm.DB, pin string) (*Session, error) {
	var s Session
	err := db.First(&s, "pin = ?", pin).Error
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
