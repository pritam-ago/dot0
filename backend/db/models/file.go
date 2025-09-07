package models

import (
	"github.com/google/uuid"
)

// File represents a file or folder in the system
type File struct {
	Base
	UserID   uuid.UUID  `gorm:"type:uuid;not null"`
	ParentID *uuid.UUID `gorm:"type:uuid"`
	Name     string     `gorm:"type:varchar(255);not null"`
	Path     string     `gorm:"type:varchar(1024);not null"`
	Size     int64      `gorm:"not null;default:0"`
	MimeType string     `gorm:"type:varchar(100)"`
	Version  int        `gorm:"not null;default:1"`

	// Relations
	User   User   `gorm:"foreignKey:UserID"`
	Parent *File  `gorm:"foreignKey:ParentID"`
	Files  []File `gorm:"foreignKey:ParentID"` // For folders
}
