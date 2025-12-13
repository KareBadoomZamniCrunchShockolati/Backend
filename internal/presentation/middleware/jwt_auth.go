package middleware

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/pkg/security"
	"strings"

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
			c.Error(exception.NewMissingUserIDException())
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(exception.NewUnauthorizedException(
				"Invalid Authorization header format. Must be 'Bearer [token]'.",
				"AUTH_HEADER_MALFORMED",
			))
			c.Abort()
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(exception.NewUnauthorizedException(
				"Invalid Authorization header format. Must be 'Bearer [token]'.",
				"AUTH_HEADER_MALFORMED",
			))
			c.Abort()
			return
		}
		
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			c.Error(exception.NewUnauthorizedException(
				"Empty token provided after 'Bearer'.",
				"TOKEN_EMPTY",
			))
			c.Abort()
			return
		}
		claims, err := m.JWTService.ValidateToken(tokenStr)
		if err != nil {
			c.Error(exception.NewUnauthorizedException(
				"Invalid or expired access token.",
				"TOKEN_INVALID_EXPIRED",
			))
			c.Abort()
			return
		}
		

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
