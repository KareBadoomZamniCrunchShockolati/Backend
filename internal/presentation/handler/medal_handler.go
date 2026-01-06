package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/exception"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MedalHandler struct {
	medalService serviceinterface.MedalServicer
}

func NewMedalHandler(medalService serviceinterface.MedalServicer) *MedalHandler {
	return &MedalHandler{medalService: medalService}
}

func (h *MedalHandler) GetUserMedals(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}
	medals, err := h.medalService.GetUserMedals(userID.(uint))
	if err != nil {
		panic(err)
	}
	Response(ctx, http.StatusOK, "", medals)
}

func (h *MedalHandler) GetSelectedMedals(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}
	medals, err := h.medalService.GetSelectedMedals(userID.(uint))
	if err != nil {
		panic(err)
	}
	Response(ctx, http.StatusOK, "", gin.H{"selected": medals})
}

func (h *MedalHandler) SelectMedals(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}
	var p dto.SelectMedalsDTO
	p = Validated[dto.SelectMedalsDTO](ctx)
	if len(p.Medals) > 3 {
		panic(exception.NewBadRequestException("TOO_MANY_MEDALS", nil))
	}
	if err := h.medalService.SelectMedals(userID.(uint), p.Medals); err != nil {
		panic(err)
	}
	Response(ctx, http.StatusOK, "", nil)
}

func (h *MedalHandler) DeleteSelectedMedals(ctx *gin.Context) {
	userID, ok := ctx.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}
	var p dto.SelectMedalsDTO
	p = Validated[dto.SelectMedalsDTO](ctx)
	if len(p.Medals) == 0 {
		panic(exception.NewBadRequestException("NO_MEDALS_PROVIDED", nil))
	}
	if err := h.medalService.DeleteSelectedMedals(userID.(uint), p.Medals); err != nil {
		panic(err)
	}
	Response(ctx, http.StatusOK, "", nil)
}
