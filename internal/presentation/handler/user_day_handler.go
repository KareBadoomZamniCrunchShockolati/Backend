package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDayHandler struct {
	userDayService serviceinterface.UserDayServicer
}

func NewUserDayHandler(userDayService serviceinterface.UserDayServicer) *UserDayHandler {
	return &UserDayHandler{
		userDayService: userDayService,
	}
}

func (h *UserDayHandler) SaveDayData(c *gin.Context) {
	type req struct {
		ChID       uint    `uri:"id"`
		Date     string  `json:"date" validate:"required,datetime=2006-01-02"`
		Note     *string `json:"note"`
		Feeling  *string `json:"feeling"`
		Progress *uint   `json:"progress"`
	}
	p := Validated[req](c)

	uID, ok := c.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}

	cID := p.ChID
	if val, exists := c.Get("challengeID"); exists {
		cID = val.(uint)
	}

	date, _ := time.Parse("2006-01-02", p.Date)
	userID := uID.(uint)

	if p.Note != nil {
		if err := h.userDayService.CreateNote(userID, cID, date, *p.Note); err != nil {
			panic(err)
		}
	}

	if p.Feeling != nil {
		val := *p.Feeling
		if val != "good" && val != "bad" {
			panic(exception.NewBadRequestException("Invalid feeling", "INVALID_FEELING", nil))
		}
		if err := h.userDayService.CreateFeeling(userID, cID, date, enum.UserFeeling(val)); err != nil {
			panic(err)
		}
	}

	if p.Progress != nil {
		if err := h.userDayService.CreateGoalProgress(userID, cID, date, *p.Progress); err != nil {
			panic(err)
		}
	}

	Response(c, 200, "Day data saved", nil)
}

func (h *UserDayHandler) UpdateDayData(ctx *gin.Context) {
    type req struct {
        ID       uint    `uri:"id"`
        Date     string  `json:"date" validate:"required,datetime=2006-01-02"`
        Note     *string `json:"note"`
        Feeling  *string `json:"feeling"`
        Progress *uint   `json:"progress"`
    }
    p := Validated[req](ctx)

    uID, ok := ctx.Get("userID")
    if !ok {
        panic(exception.NewMissingUserIDException())
    }

    cID := p.ID
    if val, exists := ctx.Get("challengeID"); exists {
        cID = val.(uint)
    }

    date, _ := time.Parse("2006-01-02", p.Date)
    userID := uID.(uint)

    if p.Note != nil {
        if err := h.userDayService.UpdateNote(userID, cID, date, *p.Note); err != nil {
            panic(err)
        }
    }
    if p.Feeling != nil {
        if err := h.userDayService.UpdateFeeling(userID, cID, date, enum.UserFeeling(*p.Feeling)); err != nil {
            panic(err)
        }
    }
    if p.Progress != nil {
        if err := h.userDayService.UpdateGoalProgress(userID, cID, date, *p.Progress); err != nil {
            panic(err)
        }
    }

    Response(ctx, 200, "Updated successfully", nil)
}

func (h *UserDayHandler) GetDayData(c *gin.Context) {
	type req struct {
		ID   uint   `uri:"id"`
		Date string `uri:"date" validate:"required,datetime=2006-01-02"`
	}
	p := Validated[req](c)

	uID, ok := c.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}

	cID := p.ID
	if val, exists := c.Get("challengeID"); exists {
		cID = val.(uint)
	}

	date, _ := time.Parse("2006-01-02", p.Date)
	userID := uID.(uint)

	note, err := h.userDayService.GetNote(userID, cID, date)
	if err != nil {
		panic(err)
	}
	feeling, err := h.userDayService.GetFeeling(userID, cID, date)
	if err != nil {
		panic(err)
	}
	progress, err := h.userDayService.GetGoalProgress(userID, cID, date)
	if err != nil {
		panic(err)
	}

	res := dto.GetDayResponse{
		Date:     p.Date,
		Note:     note,
		Feeling:  feeling,
		Progress: progress,
	}

	Response(c, 200, "", res)
}

func (h *UserDayHandler) DeleteDayData(c *gin.Context) {
	type req struct {
		ID   uint   `uri:"id"`
		Date string `uri:"date" validate:"required,datetime=2006-01-02"`
	}
	p := Validated[req](c)

	uID, ok := c.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}

	cID := p.ID
	if val, exists := c.Get("challengeID"); exists {
		cID = val.(uint)
	}

	date, _ := time.Parse("2006-01-02", p.Date)
	userID := uID.(uint)

	h.userDayService.DeleteNote(userID, cID, date)
	h.userDayService.DeleteFeeling(userID, cID, date)
	h.userDayService.DeleteGoalProgress(userID, cID, date)

	Response(c, 200, "Day data deleted successfully", nil)
}

func (h *UserDayHandler) GetGoalProgressChart(c *gin.Context) {
	type req struct {
		ID    uint   `uri:"id"`
		Start string `form:"start" validate:"required,datetime=2006-01-02"`
		End   string `form:"end" validate:"required,datetime=2006-01-02"`
	}
	p := Validated[req](c)

	uID, ok := c.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}

	cID := p.ID
	if val, exists := c.Get("challengeID"); exists {
		cID = val.(uint)
	}

	sDate, _ := time.Parse("2006-01-02", p.Start)
	eDate, _ := time.Parse("2006-01-02", p.End)

	chart, err := h.userDayService.GetGoalProgressChart(uID.(uint), cID, sDate, eDate)
	if err != nil {
		panic(err)
	}

	Response(c, 200, "", chart)
}

func (h *UserDayHandler) GetFeelingCounts(c *gin.Context) {
	type req struct {
		ID uint `uri:"id"`
	}
	p := Validated[req](c)

	uID, ok := c.Get("userID")
	if !ok {
		panic(exception.NewMissingUserIDException())
	}

	cID := p.ID
	if val, exists := c.Get("challengeID"); exists {
		cID = val.(uint)
	}

	counts, err := h.userDayService.GetTotalFeelings(uID.(uint), cID)
	if err != nil {
		panic(err)
	}

	Response(c, 200, "", counts)
}