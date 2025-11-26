package serviceinterface

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/model"
)

type ChallengeServicer interface {
	CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error)
	UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error)
	DeleteChallenge(challengeID uint, currentUserID uint) error
	GetChallengeByID(id uint, currentUserID uint) (*dto.ChallengeDetailDTO, error)
	GetAllChallenges(userID, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	StopChallenge(challengeID, currentUserID uint) error

	ListDiscoverableChallenges(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListPublicChallenges(offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCreatorID(creatorID, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCreatorUsername(username string, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCategoryID(categoryID, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByCategoryName(name string, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByParticipantCount(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesByLikeCount(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesStartingSoon(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListTopCreatorsChallenge(offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	ListChallengesJoinedByUser(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)
	SearchChallenges(query string, visibility []enum.ChallengeVisibility, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error)

	JoinPublicChallenge(userID, challengeID uint) error
	JoinPrivateChallenge(userID, challengeID uint) error
	InviteUserToChallenge(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error)
	RemoveParticipant(challengeID, removerID, participantID uint) error
	ListChallengeParticipants(challengeID uint, userID uint, offset, limit int) ([]*model.ChallengeParticipant, error)

	AcceptJoinRequest(requestID, currentUserID uint) error
	DeclineJoinRequest(requestID, currentUserID uint) error
	AcceptInvite(inviteID, currentUserID uint) error
	DeclineInvite(inviteID, currentUserID uint) error

	LeaveChallenge(userID, challengeID uint) error

	AddComment(userID uint, input *dto.AddCommentDTO) (*model.ChallengeComment, error)
	GetAllComments(challengeID, userID uint, offset, limit int) ([]*model.ChallengeComment, error)
	GetComment(commentID, userID uint) (*dto.CommentResponseDTO, error)

	GetAllCategories() ([]*model.ChallengeCategoryModel, error)

	GetRequestsSentByUser(userID uint) ([]*model.ChallengeRequest, error)
	GetInvitesSentToUser(userID uint) ([]*model.ChallengeInvite, error)
	GetInvitesSentFromChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeInvite, error)
	GetRequestsSentToChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeRequest, error)

	IsUserParticipant(challengeID, userID uint) (bool, error)
	IsChallengeCreator(challengeID, userID uint) (bool, error)
	GetChallengeParticipantCount(challengeID uint) (int, error)
	GetChallengesUserIsParticipating(userID uint, offset, limit int) ([]*model.ChallengeModel, error)
	GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error)

	LikeChallenge(userID, challengeID uint) error
	UnlikeChallenge(userID, challengeID uint) error
	IsUserLikedChallenge(userID, challengeID uint) (bool, error)
	GetChallengeLikeCount(challengeID uint) (uint, error)
}
