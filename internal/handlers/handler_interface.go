package handlers

import (
	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

type UserHandlerInterface interface {
	GetProfile(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteUser(c *gin.Context)
}
