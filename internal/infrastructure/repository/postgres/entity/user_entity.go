package entity

import (
	"time"
)

type UserEntity struct {
	ID             uint   `gorm:"type:uint;primaryKey"`
	Username       string `gorm:"unique;not null;type:varchar(50)"`
	Email          string `gorm:"unique;not null;type:varchar(100)"`
	PasswordHash   string `gorm:"not null"`
	Bio            string `gorm:"type:text"`
	ProfilePicture string `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (UserEntity) TableName() string {
	return "users"
}
