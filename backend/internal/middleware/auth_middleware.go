package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pritam-ago/dot0/backend/internal/services"
)

const (
	authHeader     = "Authorization"
	bearerSchema   = "Bearer "
	userContextKey = "user"
)

// AuthMiddleware creates a middleware for JWT authentication
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(authHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no authorization header"})
			return
		}

		if !strings.HasPrefix(header, bearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		token := header[len(bearerSchema):]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Get user from database to ensure they still exist
		user, err := authService.GetUserByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		// Store user in context
		c.Set(userContextKey, user)
		c.Next()
	}
}

// GetAuthUser retrieves the authenticated user from the context
func GetAuthUser(c *gin.Context) interface{} {
	return c.MustGet(userContextKey)
}
