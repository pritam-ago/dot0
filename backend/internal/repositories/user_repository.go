package repositories

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/pritam-ago/dot0/backend/db/models"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailTaken      = errors.New("email already taken")
	ErrInvalidProvider = errors.New("invalid auth provider")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		// Check for unique constraint violation on email
		if r.db.Where("email = ?", user.Email).First(&models.User{}).Error == nil {
			return ErrEmailTaken
		}
		return result.Error
	}
	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	return result.Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}
