package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetProfile(c *gin.Context)
	GetAllUsers(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteUser(c *gin.Context)
}
