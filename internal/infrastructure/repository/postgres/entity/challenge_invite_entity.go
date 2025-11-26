package entity

import "gorm.io/gorm"

type ChallengeInviteEntity struct {
	gorm.Model
	ChallengeID uint `gorm:"not null;index"`
	InviterID   uint `gorm:"not null;index"`
	InviteeID   uint `gorm:"not null;index"`
	Status      uint `gorm:"not null;default:1"`

	Challenge ChallengeEntity `gorm:"foreignKey:ChallengeID"`
	Inviter   UserEntity      `gorm:"foreignKey:InviterID"`
	Invitee   UserEntity      `gorm:"foreignKey:InviteeID"`
}

func (ChallengeInviteEntity) TableName() string {
	return "challenge_invites"
}
