package entity

import "time"

type UserMedalEntity struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index"`
	CategoryID uint
	Type       string
	AwardedAt  time.Time
	CreatedAt  time.Time
}

type UserCategoryProgressEntity struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"index"`
	CategoryID uint
	Count      int
	UpdatedAt  time.Time
}

type UserSelectedMedalsEntity struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"uniqueIndex"`
	Medals    string
	UpdatedAt time.Time
}
