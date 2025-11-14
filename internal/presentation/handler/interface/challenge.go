package handler

import (
	"github.com/gin-gonic/gin"
)

type ChallengeHandler interface {
	CreateChallenge(ctx *gin.Context)
	GetChallengeByID(ctx *gin.Context)
	UpdateChallenge(ctx *gin.Context)
	DeleteChallenge(ctx *gin.Context)
	GetAllChallenges(ctx *gin.Context)
	ListByCategory(ctx *gin.Context)
	ListByCreator(ctx *gin.Context)
	StopChallenge(ctx *gin.Context)
	JoinPublicChallenge(ctx *gin.Context)
	JoinPrivateChallenge(ctx *gin.Context)
	InviteUserToChallenge(ctx *gin.Context)
	RemoveParticipant(ctx *gin.Context)
	ListChallengeParticipants(ctx *gin.Context)
	AcceptJoinRequest(ctx *gin.Context)
	DeclineJoinRequest(ctx *gin.Context)
	AcceptInvite(ctx *gin.Context)
	DeclineInvite(ctx *gin.Context)
	LeaveChallenge(ctx *gin.Context)
	AddComment(ctx *gin.Context)
	GetAllComments(ctx *gin.Context)
	GetRequestsSentByUser(ctx *gin.Context)
	GetInvitesSentToUser(ctx *gin.Context)
	GetInvitesSentFromChallenge(ctx *gin.Context)
	GetRequestsSentToChallenge(ctx *gin.Context)
	GetChallengesUserIsParticipating(ctx *gin.Context)
}
