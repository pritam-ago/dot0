package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pritam-ago/dot0/backend/db/models"
	"github.com/pritam-ago/dot0/backend/internal/middleware"
	"github.com/pritam-ago/dot0/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type signupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Signup handles user registration
func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.Signup(services.SignupInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch err {
		case services.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		}
		return
	}

	c.JSON(http.StatusCreated, authResponse{
		User:  result.User,
		Token: result.Token,
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.Login(services.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		User:  result.User,
		Token: result.Token,
	})
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *gin.Context) {
	user := middleware.GetAuthUser(c).(*models.User)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GoogleAuth handles Google OAuth (stub for future implementation)
func (h *AuthHandler) GoogleAuth(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Google OAuth coming soon"})
}

// GithubAuth handles GitHub OAuth (stub for future implementation)
func (h *AuthHandler) GithubAuth(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "GitHub OAuth coming soon"})
}
