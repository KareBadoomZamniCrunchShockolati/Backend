package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	GetProfile(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteUser(c *gin.Context)
}
