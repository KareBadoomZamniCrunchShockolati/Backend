package middleware

import (
	"challenge-app/pkg/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 2. Validate the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := security.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 3. Set the user ID in the context
		if userID, ok := (*claims)["userID"]; ok {
			c.Set("userID", userID)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID not found in token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
