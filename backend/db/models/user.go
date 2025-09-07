package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	Base
	Email        string  `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash *string `gorm:"type:varchar(255)"` // Nullable for OAuth users
	AuthProvider string  `gorm:"type:varchar(50);not null;default:'local'"`

	// Relations
	UserPlan        *UserPlan        `gorm:"foreignKey:UserID"`
	PairingSessions []PairingSession `gorm:"foreignKey:UserID"`
	Devices         []Device         `gorm:"foreignKey:UserID"`
	Files           []File           `gorm:"foreignKey:UserID"`
}

// UserPlan represents a user's subscription plan
type UserPlan struct {
	Base
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	PlanType    string    `gorm:"type:varchar(50);not null"`
	Status      string    `gorm:"type:varchar(20);not null;default:'active'"`
	PurchasedAt time.Time `gorm:"not null"`
	ExpiresAt   *time.Time

	// Relations
	User User `gorm:"foreignKey:UserID"`
}
