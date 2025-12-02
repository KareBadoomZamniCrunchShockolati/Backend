package dto

import "challenge-app/internal/domain/enum"

type CreateOrUpdateDayRequest struct {
	Date     string           `json:"date" binding:"required,datetime=2006-01-02"`
	Note     *string          `json:"note,omitempty"`
	Feeling  *enum.UserFeeling `json:"feeling,omitempty"`
	Progress *uint            `json:"progress,omitempty"`
}

type GetDayResponse struct {
	Date     string           `json:"date"`
	Note     string           `json:"note,omitempty"`
	Feeling  enum.UserFeeling `json:"feeling,omitempty"`
	Progress uint             `json:"progress,omitempty"`
}

type DailyProgress struct {
	Date     string `json:"date"`
	Progress uint   `json:"progress"`
}

type DaysProgressResponse struct {
	Dates        []string `json:"dates"`
	Progress     []uint   `json:"progress"`
	MeanProgress float64  `json:"mean_progress"`
}

type FeelingCountResponse struct {
	GoodDays int `json:"good_days"`
	BadDays  int `json:"bad_days"`
}