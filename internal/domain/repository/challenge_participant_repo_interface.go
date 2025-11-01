package repository

import "challenge-app/internal/domain/model"

type ChallengeParticipantRepository interface {
	// CRUD
	AddParticipant(participant *model.ChallengeParticipant) error
	UpdateParticipantStatus(participant *model.ChallengeParticipant) error
	RemoveParticipant(challengeID, userID uint) error
	GetParticipant(challengeID, userID uint) (*model.ChallengeParticipant, error)

	// Queries
	ListParticipants(challengeID uint) ([]*model.ChallengeParticipant, error)
	ListParticipantsInUserFollowings(challengeID, userID uint) ([]*model.ChallengeParticipant, error)
	ListUserChallenges(userID uint) ([]*model.ChallengeParticipant, error) // joined challenges
}
