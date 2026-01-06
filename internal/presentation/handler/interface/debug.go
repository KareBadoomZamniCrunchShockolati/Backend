package handler

import "github.com/gin-gonic/gin"

type DebugHandler interface {
	SendTestNotification(c *gin.Context)
}
