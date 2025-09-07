package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/pritam-ago/dot0/backend/internal/handlers"
	"github.com/pritam-ago/dot0/backend/internal/middleware"
	"github.com/pritam-ago/dot0/backend/internal/services"
)

// SetupAuthRoutes configures all authentication related routes
func SetupAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	handler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.AuthMiddleware(authService)

	// Auth routes group
	auth := router.Group("/auth")
	{
		// Public routes
		auth.POST("/signup", handler.Signup)
		auth.POST("/login", handler.Login)

		// OAuth routes (stubs)
		auth.GET("/google", handler.GoogleAuth)
		auth.GET("/github", handler.GithubAuth)

		// Protected routes
		auth.GET("/me", authMiddleware, handler.Me)
	}
}
