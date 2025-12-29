package entity

import "time"

type MedalEntity struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	CategoryID  uint
	Type        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
