package model

import "time"

type MedalModel struct {
	ID          uint
	Name        string
	Description string
	CategoryID  uint
	Type        string // these 4 types: rookie, bronze, silver, gold
}

type UserMedalModel struct {
	ID         uint
	UserID     uint
	CategoryID uint
	Type       string
	AwardedAt  time.Time
}

type UserCategoryProgress struct {
	ID         uint
	UserID     uint
	CategoryID uint
	Count      int
}
