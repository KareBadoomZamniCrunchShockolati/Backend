package handler

import "github.com/gin-gonic/gin"

type ChallengeCompletionHandler interface {
	GetCompletedChallenges(ctx *gin.Context)
	GetCompletionStats(ctx *gin.Context)
	CheckChallengeCompletion(ctx *gin.Context)
}
