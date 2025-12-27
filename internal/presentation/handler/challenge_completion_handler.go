package handler

import (
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChallengeCompletionHandler struct {
	completionService serviceinterface.ChallengeCompletionServicer
}

func NewChallengeCompletionHandler(
	completionService serviceinterface.ChallengeCompletionServicer,
) *ChallengeCompletionHandler {
	return &ChallengeCompletionHandler{
		completionService: completionService,
	}
}

// completed challenges for the authenticated user
func (h *ChallengeCompletionHandler) GetCompletedChallenges(ctx *gin.Context) {
	type params struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}

	p := Validated[params](ctx)

	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	offset, limit := GetOffsetLimit(p.Page, p.PageSize, 1, 10)

	challenges, err := h.completionService.GetUserCompletedChallenges(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

// user's challenge completion statistics
func (h *ChallengeCompletionHandler) GetCompletionStats(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	stats, err := h.completionService.GetUserCompletionStats(userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", stats)
}

// checks if user completed a specific challenge
func (h *ChallengeCompletionHandler) CheckChallengeCompletion(ctx *gin.Context) {
	type uriParams struct {
		ChallengeID uint `uri:"challenge_id" validate:"required"`
	}

	p := Validated[uriParams](ctx)

	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	isCompleted, completionRate, err := h.completionService.CheckUserCompletion(
		userID.(uint),
		p.ChallengeID,
	)

	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", gin.H{
		"is_completed":    isCompleted,
		"completion_rate": completionRate,
	})
}
