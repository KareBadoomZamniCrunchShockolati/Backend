package entity

import (
	"challenge-app/internal/domain/enum"
	"time"
)

type ChallengeRequestEntity struct {
	ID          uint                   `gorm:"primaryKey;autoIncrement"`
	ChallengeID uint                   `gorm:"not null;index"`
	RequesterID uint                   `gorm:"not null;index"`
	Status      enum.RequestStatus     `gorm:"not null;type:int"`
	CreatedAt   time.Time              `gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime"`

	Challenge   *ChallengeEntity `gorm:"foreignKey:ChallengeID;constraint:OnDelete:CASCADE"`
	Requester   *UserEntity      `gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
}

func (ChallengeRequestEntity) TableName() string {
	return "challenge_requests"
}
