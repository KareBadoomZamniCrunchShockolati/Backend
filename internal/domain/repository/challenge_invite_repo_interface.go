package repository

import (
	"challenge-app/internal/domain/model"
)

type ChallengeInviteRepository interface {
	CreateInvite(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error)
	GetInvite(inviteID uint) (*model.ChallengeInvite, error)
	GetInvitesSentToUser(userID uint) ([]*model.ChallengeInvite, error)
	GetInvitesSentFromChallenge(challengeID uint) ([]*model.ChallengeInvite, error)
	UpdateInvite(invite *model.ChallengeInvite) (*model.ChallengeInvite, error)
	DeleteInvite(inviteID uint) error
}
