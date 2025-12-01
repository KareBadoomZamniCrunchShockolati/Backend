package entity

import (
	"gorm.io/gorm"
)

type CommentEntity struct {
	gorm.Model
	EntityType string `gorm:"size:20;not null;index:idx_entity"`
	EntityID   uint   `gorm:"not null;index:idx_entity"`
	UserID     uint   `gorm:"not null;index"`
	Content    string `gorm:"type:text;not null"`
	ParentID   *uint  `gorm:"index"`

	User UserEntity `gorm:"foreignKey:UserID"`
}

func (CommentEntity) TableName() string {
	return "comments"
}
