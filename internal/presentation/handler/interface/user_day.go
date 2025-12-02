package handler

import "github.com/gin-gonic/gin"

type UserDayHandler interface {
	SaveDayData(ctx *gin.Context)
	UpdateDayData(ctx *gin.Context)
	GetDayData(ctx *gin.Context)
	DeleteDayData(ctx *gin.Context)
	GetGoalProgressChart(ctx *gin.Context)
	GetFeelingCounts(ctx *gin.Context)
}