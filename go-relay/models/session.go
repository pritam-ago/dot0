package models

import (
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
