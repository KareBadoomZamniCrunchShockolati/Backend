package models

import (
	"time"
	"github.com/google/uuid"
)

type UserModel struct {
    // We use uuid.UUID for the primary key
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username     string    `gorm:"unique;not null;type:varchar(50)"`
	Email        string    `gorm:"unique;not null;type:varchar(100)"`
	PasswordHash string    `gorm:"not null"`
	Bio          string    `gorm:"type:text"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName overrides the table name
func (UserModel) TableName() string {
	return "users"
}