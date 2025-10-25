package handler

import (
	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
	Verify(c *gin.Context)
	ResendVerification(c *gin.Context)
}