package entity

import "gorm.io/gorm"

type FollowEntity struct {
	gorm.Model
	FollowerID  uint `gorm:"not null; index"`
	FollowingID uint `gorm:"not null; index"`

	Follower  UserEntity `gorm:"foreignKey:FollowerID;references:ID"`
	Following UserEntity `gorm:"foreignKey:FollowingID;references:ID"`
}

func (FollowEntity) TableName() string {
	return "follows"
}

//no repetative follows

func (FollowEntity) TableConsstraints() []string {
	return []string{
		"UNIQUE (follower_id, following_id)",
	}
}
