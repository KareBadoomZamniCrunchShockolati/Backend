package middleware

import (
	"strings"

	"challenge-app/internal/domain/exception"
	"challenge-app/pkg/security"

	"github.com/gin-gonic/gin"
)

type JWTMiddleware struct {
	JWTService security.JWTService
}

func NewJWTMiddleware(jwtSvc security.JWTService) *JWTMiddleware {
	return &JWTMiddleware{
		JWTService: jwtSvc,
	}
}

func (m *JWTMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			panic(exception.NewUnauthorizedException(
				"AUTH_MISSING_ID",
				map[string]any{
					"header": "Authorization",
				},
			))
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			panic(exception.NewUnauthorizedException(
				"AUTH_HEADER_MALFORMED",
				map[string]any{
					"expected": "Bearer <token>",
				},
			))
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenStr == "" {
			panic(exception.NewUnauthorizedException(
				"TOKEN_EMPTY",
				nil,
			))
		}
		claims, err := m.JWTService.ValidateToken(tokenStr)
		if err != nil {
			panic(exception.NewUnauthorizedException(
				"TOKEN_INVALID_EXPIRED",
				map[string]any{
					"reason": err.Error(),
				},
			))
		}

		if claims.UserID == 0 {
			panic(exception.NewUnauthorizedException(
				"TOKEN_INVALID_PAYLOAD",
				nil,
			))
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
