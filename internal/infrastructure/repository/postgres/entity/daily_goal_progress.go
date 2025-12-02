package entity
import (
	"time"
	"gorm.io/gorm"
)	

type UserDailyGoalProgressEntity struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index:idx_progress_user_date"`
	ChallengeID uint      `gorm:"not null;index:idx_progress_user_date"`
	Date        time.Time `gorm:"not null;index:idx_progress_user_date;type:date"`

	Progress uint `gorm:"not null"` 

	User      UserEntity      `gorm:"foreignKey:UserID"`
	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
}

func (UserDailyGoalProgressEntity) TableName() string {
	return "user_daily_goal_progress"
}