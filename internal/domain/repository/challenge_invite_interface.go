package repository

import "challenge-app/internal/domain/model"

type ChallengeInviteRepository interface {
	CreateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error)
	GetInvite(inviteID uint) (*model.ChallengeInvite, error)
	GetUserInvites(userID uint) ([]*model.ChallengeInvite, error)
	GetChallengeInvites(challengeID uint) ([]*model.ChallengeInvite, error)
	UpdateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error)
	DeleteInvite(inviteID uint) error
}
