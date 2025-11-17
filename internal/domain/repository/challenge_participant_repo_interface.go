package repository

import "challenge-app/internal/domain/model"

type ChallengeParticipantRepository interface {
	CreateParticipant(participant *model.ChallengeParticipant) (*model.ChallengeParticipant, error)
	GetParticipant(challengeID, userID uint) (*model.ChallengeParticipant, error)
	GetParticipantsByChallenge(challengeID uint) ([]*model.ChallengeParticipant, error)
	UpdateParticipant(participant *model.ChallengeParticipant) (*model.ChallengeParticipant, error)
	DeleteParticipant(challengeID, userID uint) error
	GetParticipantCount(challengeID uint) (int, error)

	// Check if user is participant
	IsUserParticipant(challengeID, userID uint) (bool, error)
}
