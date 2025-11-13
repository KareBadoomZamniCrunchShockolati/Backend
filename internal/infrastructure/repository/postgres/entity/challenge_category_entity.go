package entity

import "gorm.io/gorm"

type ChallengeCategoryEntity struct {
	gorm.Model 
	Name string `gorm:"size:100;not null;uniqueIndex"` 
	Description string `gorm:"size:255"` 
}

func (ChallengeCategoryEntity) TableName() string {
	return "challenge_categories"
}