package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayeed1999/share-a-ride/internal/domain/entity"
)

// RequireRole middleware checks if the user has the required role
func RequireRole(roles ...entity.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by auth middleware)
		roleInterface, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userRole := roleInterface.(entity.Role)
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IsAdmin checks if the user has admin role
func IsAdmin() gin.HandlerFunc {
	return RequireRole(entity.RoleAdmin)
}

// IsUser checks if the user has user role
func IsUser() gin.HandlerFunc {
	return RequireRole(entity.RoleUser, entity.RoleAdmin)
}
