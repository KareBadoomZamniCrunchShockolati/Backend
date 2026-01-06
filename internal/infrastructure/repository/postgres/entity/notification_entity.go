package entity

import "time"

type NotificationEntity struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"index;not null"`
	Type     string `gorm:"type:varchar(64);not null"`
	TitleKey string `gorm:"type:varchar(255);not null;default:''"`
	BodyKey  string `gorm:"type:varchar(255);not null;default:''"`
	Data      []byte `gorm:"type:jsonb"`
	CreatedAt time.Time `gorm:"index;not null"`
	ReadAt    *time.Time `gorm:"index"`
}


func (NotificationEntity) TableName() string { return "notifications" }
