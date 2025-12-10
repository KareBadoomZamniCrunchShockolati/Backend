package handler

import (
	"challenge-app/internal/application/dto"
	serviceinterface "challenge-app/internal/application/service/interface"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	challengeService serviceinterface.ChallengeServicer
}

func NewChallengeHandler(
	challengeService serviceinterface.ChallengeServicer,
) *ChallengeHandler {
	return &ChallengeHandler{
		challengeService: challengeService,
	}
}

func (h *ChallengeHandler) CreateChallenge(ctx *gin.Context) {
	type createChallengeParams struct {
		Title           string `json:"title" validate:"required,min=5,max=150"`
		Description     string `json:"description" validate:"required,min=10,max=280"`
		CategoryID      uint   `json:"category_id" validate:"required"`
		MaxParticipants uint   `json:"max_participants" validate:"omitempty,min=0"`
		Visibility      string `json:"visibility" validate:"required,oneof=public private invite"`
		Location        string `json:"location"`
		Goal            *int   `json:"goal"`
		Rule            string `json:"rule" validate:"required"`
		CommentsEnabled bool   `json:"comments_enabled"`
		StartTime       string `json:"start_time" validate:"required"`
		EndTime         string `json:"end_time" validate:"required"`
		Timezone        string `json:"timezone"`
		ImageURL        string `json:"image_url"`
	}
	params := Validated[createChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	startTime, err := time.Parse(time.RFC3339, params.StartTime)
	if err != nil {
		panic(exception.NewBadRequestException("Invalid start time format", "INVALID_TIME_FORMAT", nil))
	}
	endTime, err := time.Parse(time.RFC3339, params.EndTime)
	if err != nil {
		panic(exception.NewBadRequestException("Invalid end time format", "INVALID_TIME_FORMAT", nil))
	}

	createChallengeDTO := &dto.CreateChallengeDTO{
		CreatorID:       userID.(uint),
		Title:           params.Title,
		Description:     params.Description,
		CategoryID:      params.CategoryID,
		MaxParticipants: params.MaxParticipants,
		Visibility:      enum.ChallengeVisibility(params.Visibility),
		Location:        params.Location,
		Goal:            params.Goal,
		Rule:            params.Rule,
		CommentsEnabled: params.CommentsEnabled,
		StartTime:       startTime,
		EndTime:         endTime,
		Timezone:        params.Timezone,
		ImageURL:        params.ImageURL,
	}
	challenge, err := h.challengeService.CreateChallenge(createChallengeDTO)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Challenge created successfully", challenge)
}

func (h *ChallengeHandler) GetChallengeByID(ctx *gin.Context) {
	type getChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[getChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	challenge, err := h.challengeService.GetChallengeByID(params.ID, userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenge)
}

func (h *ChallengeHandler) UpdateChallenge(ctx *gin.Context) {
	type updateChallengeParams struct {
		ID              uint    `uri:"id" validate:"required"`
		Title           *string `json:"title"`
		Description     *string `json:"description"`
		CategoryID      *uint   `json:"category_id"`
		MaxParticipants *uint   `json:"max_participants"`
		Visibility      *uint   `json:"visibility"`
		Location        *string `json:"location"`
		Goal            *int    `json:"goal"`
		Rule            *string `json:"rule"`
		CommentsEnabled *bool   `json:"comments_enabled"`
		IsStopped       *bool   `json:"is_stopped"`
		EndTime         *string `json:"end_time"`
		ImageURL        *string `json:"image_url"`
		StartTime       *string `json:"start_time"`
		Timezone        *string `json:"timezone"`
	}
	params := Validated[updateChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	var startTime *time.Time
	var endTime *time.Time
	if params.StartTime != nil {
		parsedTime, err := time.Parse(time.RFC3339, *params.StartTime)
		if err != nil {
			panic(exception.NewBadRequestException("Invalid start time format", "INVALID_TIME_FORMAT", nil))
		}
		startTime = &parsedTime
	}
	if params.EndTime != nil {
		parsedTime, err := time.Parse(time.RFC3339, *params.EndTime)
		if err != nil {
			panic(exception.NewBadRequestException("Invalid end time format", "INVALID_TIME_FORMAT", nil))
		}
		endTime = &parsedTime
	}
	var visibility *enum.ChallengeVisibility
	if params.Visibility != nil {
		v := enum.ChallengeVisibility(*params.Visibility)
		visibility = &v
	}
	updateChallengeDTO := &dto.UpdateChallengeDTO{
		Title:           params.Title,
		Description:     params.Description,
		CategoryID:      params.CategoryID,
		MaxParticipants: params.MaxParticipants,
		Visibility:      visibility,
		Location:        params.Location,
		Goal:            params.Goal,
		Rule:            params.Rule,
		CommentsEnabled: params.CommentsEnabled,
		IsStopped:       params.IsStopped,
		EndTime:         endTime,
		ImageURL:        params.ImageURL,
		StartTime:       startTime,
		Timezone:        params.Timezone,
	}

	updatedChallenge, err := h.challengeService.UpdateChallenge(params.ID, userID.(uint), updateChallengeDTO)
	if err != nil {
		panic(err)
	}
	Response(ctx, http.StatusOK, "Challenge updated successfully", updatedChallenge)
}

func (h *ChallengeHandler) DeleteChallenge(ctx *gin.Context) {
	type deleteChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[deleteChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.DeleteChallenge(params.ID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Challenge deleted successfully", nil)
}

func (h *ChallengeHandler) ListDiscoverableChallenges(ctx *gin.Context) {
	type discoverableParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[discoverableParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListDiscoverableChallenges(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListPublicChallenges(ctx *gin.Context) {
	type publicParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[publicParams](ctx)
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListPublicChallenges(offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByCreator(ctx *gin.Context) {
	type listByCreatorParams struct {
		UserID   uint `uri:"user_id" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	params := Validated[listByCreatorParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByCreatorID(params.UserID, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByCreatorUsername(ctx *gin.Context) {
	type listByCreatorUsernameParams struct {
		Username string `uri:"username" validate:"required"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	params := Validated[listByCreatorUsernameParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByCreatorUsername(params.Username, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByCategory(ctx *gin.Context) {
	type listByCategoryParams struct {
		CategoryID uint `uri:"category_id" validate:"required"`
		Page       int  `form:"page"`
		PageSize   int  `form:"pageSize"`
	}
	params := Validated[listByCategoryParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByCategoryID(params.CategoryID, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByCategoryName(ctx *gin.Context) {
	type listByCategoryNameParams struct {
		CategoryName string `uri:"category_name" validate:"required"`
		Page         int    `form:"page"`
		PageSize     int    `form:"pageSize"`
	}
	params := Validated[listByCategoryNameParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByCategoryName(params.CategoryName, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByParticipantCount(ctx *gin.Context) {
	type participantCountParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[participantCountParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByParticipantCount(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListByLikeCount(ctx *gin.Context) {
	type likeCountParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[likeCountParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesByLikeCount(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListChallengesStartingSoon(ctx *gin.Context) {
	type startingSoonParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[startingSoonParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesStartingSoon(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListTopCreatorsChallenge(ctx *gin.Context) {
	type topCreatorsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[topCreatorsParams](ctx)
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListTopCreatorsChallenge(offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) ListTopCreators(ctx *gin.Context) {
	type topCreatorsParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[topCreatorsParams](ctx)
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 20)
	creators, err := h.challengeService.ListTopCreators(offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", creators)
}

func (h *ChallengeHandler) GetChallengesUserIsParticipating(ctx *gin.Context) {
	type getParticipatingChallengesParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	params := Validated[getParticipatingChallengesParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.ListChallengesJoinedByUser(userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) SearchChallenges(ctx *gin.Context) {
	type searchParams struct {
		Query    string `form:"query" validate:"required,min=1"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	params := Validated[searchParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	}
	var visibility []enum.ChallengeVisibility
	// TODO: implement visibility parsing

	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	challenges, err := h.challengeService.SearchChallenges(params.Query, visibility, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) SearchChallengesUserIsParticipating(ctx *gin.Context) {
	type searchParams struct {
		Query    string `form:"query" validate:"required,min=1"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	params := Validated[searchParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)

	challenges, err := h.challengeService.SearchChallengesUserIsParticipating(
		userID.(uint),
		params.Query,
		offset,
		limit,
	)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", challenges)
}

func (h *ChallengeHandler) StopChallenge(ctx *gin.Context) {
	type stopChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[stopChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.StopChallenge(params.ID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Challenge stopped successfully", nil)
}

func (h *ChallengeHandler) JoinPublicChallenge(ctx *gin.Context) {
	type joinPublicChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[joinPublicChallengeParams](ctx)

	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.JoinPublicChallenge(userID.(uint), params.ID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Successfully joined challenge", nil)
}

func (h *ChallengeHandler) JoinPrivateChallenge(ctx *gin.Context) {
	type joinPrivateChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[joinPrivateChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.JoinPrivateChallenge(userID.(uint), params.ID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Join request sent successfully", nil)
}

func (h *ChallengeHandler) InviteUserToChallenge(ctx *gin.Context) {
	type inviteUserParams struct {
		ID        uint `uri:"id" validate:"required"`
		InviteeID uint `json:"invitee_id" validate:"required"`
	}
	params := Validated[inviteUserParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	createdInvite, err := h.challengeService.InviteUserToChallenge(userID.(uint), params.ID, params.InviteeID)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "User invited successfully", createdInvite)
}

func (h *ChallengeHandler) RemoveParticipant(ctx *gin.Context) {
	type removeParticipantParams struct {
		ID            uint `uri:"id" validate:"required"`
		ParticipantID uint `uri:"participant_id" validate:"required"`
	}
	params := Validated[removeParticipantParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.RemoveParticipant(params.ID, userID.(uint), params.ParticipantID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Participant removed successfully", nil)
}

func (h *ChallengeHandler) ListChallengeParticipants(ctx *gin.Context) {
	type listParticipantsParams struct {
		ID       uint `uri:"id" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	params := Validated[listParticipantsParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	participants, err := h.challengeService.ListChallengeParticipants(params.ID, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", participants)
}

func (h *ChallengeHandler) AcceptJoinRequest(ctx *gin.Context) {
	type acceptJoinRequestParams struct {
		RequestID uint `uri:"request_id" validate:"required"`
	}
	params := Validated[acceptJoinRequestParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.AcceptJoinRequest(params.RequestID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Join request accepted successfully", nil)
}

func (h *ChallengeHandler) DeclineJoinRequest(ctx *gin.Context) {
	type declineJoinRequestParams struct {
		RequestID uint `uri:"request_id" validate:"required"`
	}
	params := Validated[declineJoinRequestParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.DeclineJoinRequest(params.RequestID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Join request declined successfully", nil)
}

func (h *ChallengeHandler) AcceptInvite(ctx *gin.Context) {
	type acceptInviteParams struct {
		InviteID uint `uri:"invite_id" validate:"required"`
	}
	params := Validated[acceptInviteParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.AcceptInvite(params.InviteID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Invite accepted successfully", nil)
}

func (h *ChallengeHandler) DeclineInvite(ctx *gin.Context) {
	type declineInviteParams struct {
		InviteID uint `uri:"invite_id" validate:"required"`
	}
	params := Validated[declineInviteParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.DeclineInvite(params.InviteID, userID.(uint)); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Invite declined successfully", nil)
}

func (h *ChallengeHandler) LeaveChallenge(ctx *gin.Context) {
	type leaveChallengeParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[leaveChallengeParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	if err := h.challengeService.LeaveChallenge(userID.(uint), params.ID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Successfully left challenge", nil)
}

func (h *ChallengeHandler) AddComment(ctx *gin.Context) {
	type addCommentParams struct {
		ID      uint   `uri:"id" validate:"required"`
		Content string `json:"content" validate:"required,min=1,max=1000"`
	}
	params := Validated[addCommentParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	addCommentDTO := &dto.CommentRequestDTO{
		EntityType: "challenge",
		EntityID:   params.ID,
		Content:    params.Content,
		ParentID:   nil,
	}
	comment, err := h.challengeService.AddComment(userID.(uint), addCommentDTO)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Comment added successfully", comment)
}

func (h *ChallengeHandler) GetAllComments(ctx *gin.Context) {
	type getAllCommentsParams struct {
		ID       uint `uri:"id" validate:"required"`
		Page     int  `form:"page"`
		PageSize int  `form:"pageSize"`
	}
	params := Validated[getAllCommentsParams](ctx)

	userID, exists := ctx.Get("userID")
	if !exists {
		userID = uint(0)
	}

	offset, limit := GetOffsetLimit(params.Page, params.PageSize, 1, 10)
	comments, err := h.challengeService.GetAllComments(params.ID, userID.(uint), offset, limit)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", comments)
}

func (h *ChallengeHandler) GetComment(ctx *gin.Context) {
	type getCommentParams struct {
		CommentID uint `uri:"comment_id" validate:"required"`
	}
	params := Validated[getCommentParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	comment, err := h.challengeService.GetComment(params.CommentID, userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", comment)
}

func (h *ChallengeHandler) GetRequestsSentByUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	requests, err := h.challengeService.GetRequestsSentByUser(userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", requests)
}

func (h *ChallengeHandler) GetInvitesSentToUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	invites, err := h.challengeService.GetInvitesSentToUser(userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", invites)
}

func (h *ChallengeHandler) GetInvitesSentFromChallenge(ctx *gin.Context) {
	type getInvitesParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[getInvitesParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	invites, err := h.challengeService.GetInvitesSentFromChallenge(params.ID, userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", invites)
}

func (h *ChallengeHandler) GetRequestsSentToChallenge(ctx *gin.Context) {
	type getRequestsParams struct {
		ID uint `uri:"id" validate:"required"`
	}
	params := Validated[getRequestsParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	requests, err := h.challengeService.GetRequestsSentToChallenge(params.ID, userID.(uint))
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", requests)
}

func (h *ChallengeHandler) GetAllCategories(ctx *gin.Context) {
	categories, err := h.challengeService.GetAllCategories()
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", categories)
}

func (h *ChallengeHandler) GetMutualFollowersInChallenge(ctx *gin.Context) {
	type getMutualFollowersParams struct {
		ChallengeID uint `uri:"id" validate:"required"`
	}
	params := Validated[getMutualFollowersParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	mutualFollowers, err := h.challengeService.GetMutualFollowersInChallenge(userID.(uint), params.ChallengeID)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", mutualFollowers)
}

func (h *ChallengeHandler) LikeChallenge(ctx *gin.Context) {
	type likeChallengeParams struct {
		ChallengeID uint `uri:"id" validate:"required"`
	}
	params := Validated[likeChallengeParams](ctx)
	userIDInterface, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {
		panic(exception.NewInternalServerException(
			exception.ErrorTypeContextCastFail,
			"Invalid user ID type",
			nil,
		))
	}

	if err := h.challengeService.LikeChallenge(userID, params.ChallengeID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Challenge liked successfully", nil)
}

func (h *ChallengeHandler) UnlikeChallenge(ctx *gin.Context) {
	type unlikeChallengeParams struct {
		ChallengeID uint `uri:"id" validate:"required"`
	}
	params := Validated[unlikeChallengeParams](ctx)
	userIDInterface, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	userID, ok := userIDInterface.(uint)
	if !ok || userID == 0 {
		panic(exception.NewInternalServerException(
			exception.ErrorTypeContextCastFail,
			"Invalid user ID type",
			nil,
		))
	}

	if err := h.challengeService.UnlikeChallenge(userID, params.ChallengeID); err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "Challenge unliked successfully", nil)
}

func (h *ChallengeHandler) GetChallengeLikeCount(ctx *gin.Context) {
	type getLikeCountParams struct {
		ChallengeID uint `uri:"id" validate:"required"`
	}
	params := Validated[getLikeCountParams](ctx)

	count, err := h.challengeService.GetChallengeLikeCount(params.ChallengeID)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", gin.H{"like_count": count})
}

func (h *ChallengeHandler) IsUserLikedChallenge(ctx *gin.Context) {
	type isLikedParams struct {
		ChallengeID uint `uri:"id" validate:"required"`
	}
	params := Validated[isLikedParams](ctx)
	userID, exists := ctx.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}

	isLiked, err := h.challengeService.IsUserLikedChallenge(userID.(uint), params.ChallengeID)
	if err != nil {
		panic(err)
	}

	Response(ctx, http.StatusOK, "", gin.H{"is_liked": isLiked})
}
