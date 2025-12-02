package entity
import (
	"time"
	"gorm.io/gorm"
)

type UserDailyNoteEntity struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index:idx_note_user_date"`
	ChallengeID uint      `gorm:"not null;index:idx_note_user_date"`
	Date        time.Time `gorm:"not null;index:idx_note_user_date;type:date"`

	Note string `gorm:"type:text"` 

	User      UserEntity      `gorm:"foreignKey:UserID"`
	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
}

func (UserDailyNoteEntity) TableName() string {
	return "user_daily_notes"
}