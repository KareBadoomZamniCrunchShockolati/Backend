// internal/application/service/interface/challenge.go
package serviceinterface

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/model"
)

type ChallengeServicer interface {
	CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error)
	UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error)
	DeleteChallenge(challengeID uint, currentUserID uint) error
	GetChallengeByID(id uint) (*model.ChallengeModel, error)
	GetAllChallenges() ([]*model.ChallengeModel, error)
	ListByCategory(categoryID uint, offset, limit int) ([]*model.ChallengeModel, error)
	ListByCreator(userID uint, offset, limit int) ([]*model.ChallengeModel, error)

	StopChallenge(challengeID, currentUserID uint) error

	JoinPublicChallenge(userID, challengeID uint) error
	JoinPrivateChallenge(userID, challengeID uint) error
	InviteUserToChallenge(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error) 
	RemoveParticipant(challengeID, removerID, participantID uint) error
	ListChallengeParticipants(challengeID uint, userID uint) ([]*model.ChallengeParticipant, error)

	AcceptJoinRequest(requestID, currentUserID uint) error
	DeclineJoinRequest(requestID, currentUserID uint) error
	AcceptInvite(inviteID, currentUserID uint) error
	DeclineInvite(inviteID, currentUserID uint) error

	LeaveChallenge(userID, challengeID uint) error

	AddComment(userID uint, input *dto.AddCommentDTO) (*model.ChallengeComment, error)
	GetAllComments(challengeID uint, offset, limit int) ([]*model.ChallengeComment, error)

	GetRequestsSentByUser(userID uint) ([]*model.ChallengeRequest, error)
	GetInvitesSentToUser(userID uint) ([]*model.ChallengeInvite, error)
	GetInvitesSentFromChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeInvite, error)
	GetRequestsSentToChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeRequest, error)

	IsUserParticipant(challengeID, userID uint) (bool, error)
	IsChallengeCreator(challengeID, userID uint) (bool, error)
	GetChallengeParticipantCount(challengeID uint) (int, error)
	GetChallengesUserIsParticipating(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
}
