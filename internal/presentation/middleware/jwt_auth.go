package middleware

import (
	"strings"
	"challenge-app/pkg/security"
	"github.com/gin-gonic/gin"
)


type JWTMiddleware struct {
	JWTService security.JWTService
}

func NewJWTMiddleware(jwtSvc security.JWTService) *JWTMiddleware {
	return &JWTMiddleware{JWTService: jwtSvc}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.JWTService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
