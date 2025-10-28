package middleware

import (
	"github.com/gin-gonic/gin"
)

type ErrorMiddleware interface {
	PanicRecovery() gin.HandlerFunc
	APIErrorTranslator() gin.HandlerFunc
}
