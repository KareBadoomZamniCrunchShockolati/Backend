package repository

import "challenge-app/internal/domain/model"

type ChallengeRequestRepository interface {
	CreateJoinRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error)
	GetJoinRequest(requestID uint) (*model.ChallengeRequest, error)
	GetChallengeRequests(challengeID uint) ([]*model.ChallengeRequest, error)
	GetUserJoinRequests(userID uint) ([]*model.ChallengeRequest, error)
	UpdateJoinRequest(request *model.ChallengeRequest) (*model.ChallengeRequest, error)
	DeleteJoinRequest(requestID uint) error
}
