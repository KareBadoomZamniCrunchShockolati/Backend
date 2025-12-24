package middleware

import "github.com/gin-gonic/gin"

type ErrorMiddleware interface {
	PanicRecovery() gin.HandlerFunc
	ErrorHandler() gin.HandlerFunc
}

type LocalizationMiddleware interface {
	Handle() gin.HandlerFunc
}
