package models

import (
	"time"

	"github.com/google/uuid"
)

// PairingSession represents a device pairing session
type PairingSession struct {
	Base
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	PCID      string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(50);index;not null"`
	Status    string    `gorm:"type:varchar(20);not null;default:'pending'"`
	ExpiresAt time.Time `gorm:"not null"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// Device represents a user's device
type Device struct {
	Base
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	DeviceType string    `gorm:"type:varchar(20);not null"` // pc, mobile
	OS         string    `gorm:"type:varchar(50);not null"`
	LastSeenAt time.Time `gorm:"not null"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}
