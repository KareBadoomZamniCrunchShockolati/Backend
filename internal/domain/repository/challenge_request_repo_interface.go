package repository

import (
	"challenge-app/internal/domain/model"
)

type ChallengeJoinRequestRepository interface {
	CreateRequest(requesterID, challengeID uint) (*model.ChallengeRequest, error)
	GetRequest(requestID uint) (*model.ChallengeRequest, error)
	GetRequestsSentByUser(userID uint) ([]*model.ChallengeRequest, error)
	GetRequestsSentToChallenge(challengeID uint) ([]*model.ChallengeRequest, error)
	UpdateRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error)
	DeleteRequest(requestID uint) error
}