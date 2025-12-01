package entity

import (
	"challenge-app/internal/domain/enum"
	"gorm.io/gorm"
)

type ChallengeParticipationRequestEntity struct {
	gorm.Model
	ChallengeID uint              `gorm:"not null;index"`
	FromUserID  uint              `gorm:"not null;index"` // inviter or requester
	ToUserID    uint              `gorm:"not null;index"` // invitee or challenge owner (for requests)
	Type        enum.RequestType  `gorm:"not null;type:varchar(20)"`
	Status      enum.RequestStatus `gorm:"not null;type:varchar(20);default:'pending'"`

	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	FromUser  UserEntity       `gorm:"foreignKey:FromUserID"`
	ToUser    UserEntity       `gorm:"foreignKey:ToUserID"`
}

func (ChallengeParticipationRequestEntity) TableName() string {
	return "challenge_participation_requests"
}