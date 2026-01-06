package handler

import "github.com/gin-gonic/gin"

type WSNotificationHandler interface {
	Connect(c *gin.Context)
}
