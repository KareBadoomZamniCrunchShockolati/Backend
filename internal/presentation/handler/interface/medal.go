package handler

import "github.com/gin-gonic/gin"

type MedalHandler interface {
	GetUserMedals(c *gin.Context)
	GetSelectedMedals(c *gin.Context)
	SelectMedals(c *gin.Context)
	DeleteSelectedMedals(c *gin.Context)
}
