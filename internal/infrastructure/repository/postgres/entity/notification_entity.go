package entity

import "time"

type NotificationEntity struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"index:idx_notif_user_created,priority:1;not null"`
	Type      string     `gorm:"type:text;not null"`
	Title     string     `gorm:"type:text;not null"`
	Body      string     `gorm:"type:text;not null"`
	Data      []byte     `gorm:"type:jsonb"` 
	CreatedAt time.Time  `gorm:"index:idx_notif_user_created,priority:2,sort:desc;not null"`
	ReadAt    *time.Time `gorm:"index"`
}

func (NotificationEntity) TableName() string { return "notifications" }
