package middleware

import (
	"github.com/gin-gonic/gin"
)

type JWTMiddleware interface {
	Handler() gin.HandlerFunc
	OptionalHandler() gin.HandlerFunc
}
