package entity

import "gorm.io/gorm"

type UserEntity struct {
    gorm.Model
    Username       string `gorm:"unique;not null;type:varchar(50)"`
    Email          string `gorm:"unique;not null;type:varchar(100)"`
    PasswordHash   string `gorm:"not null"`
    Bio            string `gorm:"type:text"`
    ProfilePicture string `gorm:"type:text"`
}

func (UserEntity) TableName() string {
    return "users"
}
