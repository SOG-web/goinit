package middleware

import (
	"net/http"
	"strings"

	"github.com/SOG-web/goinit/gin/api/common/dto"
	"github.com/SOG-web/goinit/gin/internal/lib/jwt"
	"github.com/gin-gonic/gin"
)

// RequireAuth ensures a user is signed in using JWT token
func RequireAuth(jwtService jwt.JWTServiceInterface) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.AuthErrorResponse{
				Error:      "Authorization header is required",
				Success:    false,
				StatusCode: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.AuthErrorResponse{
				Error:      "Invalid authorization header format",
				Success:    false,
				StatusCode: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.AuthErrorResponse{
				Error:      "Invalid or expired token",
				Success:    false,
				StatusCode: http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		// c.Set("is_staff", claims.IsStaff)
		// c.Set("is_superuser", claims.IsSuperuser)

		c.Next()
	})
}

// RequireStaff middleware ensures user has staff privileges
func RequireAdmin() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, dto.AuthErrorResponse{
				Error:      "Admin privileges required",
				Success:    false,
				StatusCode: http.StatusForbidden,
			})
			c.Abort()
			return
		}

		c.Next()
	})
}

// RequireRole checks that the current user has the given role (legacy compatibility)
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		switch role {
		case "admin":
			if !isAdmin.(bool) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		case "superuser":
			isSuperuser, exists := c.Get("is_superuser")
			if !exists || !isSuperuser.(bool) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "unknown role"})
			return
		}

		c.Next()
	}
}
