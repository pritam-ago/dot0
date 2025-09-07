package services

import (
	"errors"

	"github.com/google/uuid"

	"github.com/pritam-ago/dot0/backend/config"
	"github.com/pritam-ago/dot0/backend/db/models"
	"github.com/pritam-ago/dot0/backend/internal/repositories"
	"github.com/pritam-ago/dot0/backend/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

type SignupInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResponse struct {
	User  *models.User
	Token string
}

func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// Signup creates a new user account
func (s *AuthService) Signup(input SignupInput) (*AuthResponse, error) {
	// Check if user exists
	if _, err := s.userRepo.FindByEmail(input.Email); err == nil {
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        input.Email,
		PasswordHash: &hashedPassword,
		AuthProvider: "email",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiryPeriod)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user has a password (might be OAuth user)
	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !utils.CheckPassword(input.Password, *user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, s.config.JWTSecret, s.config.JWTExpiryPeriod)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(token string) (*utils.Claims, error) {
	return utils.ValidateToken(token, s.config.JWTSecret)
}
