package handler

import (
	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	Verify(c *gin.Context)
	ResendVerification(c *gin.Context)
}
