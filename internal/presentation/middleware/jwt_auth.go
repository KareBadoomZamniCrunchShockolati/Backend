package middleware

import (
	"strings"
	"challenge-app/pkg/security"
	"github.com/gin-gonic/gin"
	"challenge-app/pkg/errs"
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
			panic(&errs.UnAuthorizedError{
				MessageValue: "Missing Authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.JWTService.ValidateToken(tokenStr)
		if err != nil {
			panic(&errs.UnAuthorizedError{
				MessageValue: "Empty or malformed token",
			})

		}
		

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
