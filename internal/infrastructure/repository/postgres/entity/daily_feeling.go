package entity

import (
	"challenge-app/internal/domain/enum"
	"time"

	"gorm.io/gorm"
)

type UserDailyFeelingEntity struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index:idx_feeling_user_date"`
	ChallengeID uint      `gorm:"not null;index:idx_feeling_user_date"`
	Date        time.Time `gorm:"not null;index:idx_feeling_user_date;type:date"`

	Feeling enum.UserFeeling `gorm:"type:varchar(10);not null"` 

	User      UserEntity      `gorm:"foreignKey:UserID"`
	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
}

func (UserDailyFeelingEntity) TableName() string {
	return "user_daily_feelings"
}
