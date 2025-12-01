package service

import (
	"fmt"
	"time"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
)

type ChallengeService struct {
	challengeRepo            repository.ChallengeRepository
	participantRepo          repository.ChallengeParticipantRepository
	commentRepo              repository.ChallengeCommentRepository
	inviteRepo               repository.ChallengeInviteRepository
	joinRequestRepo          repository.ChallengeJoinRequestRepository
	categoryRepo             repository.CategoryRepository
	userRepo                 repository.UserRepository
	followRepo               repository.FollowRepository
	likeRepo                 repository.LikeRepository
}

func NewChallengeService(
	challengeRepo repository.ChallengeRepository,
	participantRepo repository.ChallengeParticipantRepository,
	commentRepo repository.ChallengeCommentRepository,
	inviteRepo repository.ChallengeInviteRepository,
	joinRequestRepo repository.ChallengeJoinRequestRepository,
	categoryRepo repository.CategoryRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	likeRepo repository.LikeRepository,
) *ChallengeService {
	return &ChallengeService{
		challengeRepo:   challengeRepo,
		participantRepo: participantRepo,
		commentRepo:     commentRepo,
		inviteRepo:      inviteRepo,
		joinRequestRepo: joinRequestRepo,
		categoryRepo:    categoryRepo,
		userRepo:        userRepo,
		followRepo:      followRepo,
		likeRepo:        likeRepo,
	}
}

// Challenge CRUD operations
func (s *ChallengeService) CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "CHALLENGE_CREATE_BAD_INPUT", nil)
	}
	category, err := s.categoryRepo.GetCategoryByID(input.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, exception.NewNotFoundException("Category", fmt.Sprintf("%d", input.CategoryID), "CATEGORY_NOT_FOUND")
	}
	user, err := s.userRepo.GetUserByID(input.CreatorID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", input.CreatorID), "USER_NOT_FOUND")
	}

	challenge := &model.ChallengeModel{
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      input.CategoryID,
		CreatorID:       input.CreatorID,
		MaxParticipants: input.MaxParticipants,
		Visibility:      input.Visibility,
		Rule:            input.Rule,
		StartTime:       input.StartTime,
		EndTime:         &input.EndTime,
		Timezone:        input.Timezone,
		ImageURL:        input.ImageURL,
		IsStopped:       false,
		CommentsEnabled: input.CommentsEnabled,
	}
	created, err := s.challengeRepo.CreateChallenge(challenge)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *ChallengeService) UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(id, currentUserID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	if challenge.CreatorID != currentUserID {
		return nil, exception.NewUnauthorizedException("Only the creator can update this challenge", "CHALLENGE_UPDATE_FORBIDDEN")
	}
	if input.StartTime != nil {
		if challenge.StartTime.Before(time.Now()) {
			return nil, exception.NewBadRequestException("Cannot change StartTime, challenge already started", "CHALLENGE_ALREADY_STARTED", nil)
		}
		challenge.StartTime = *input.StartTime
	}

	if input.Title != nil {
		challenge.Title = *input.Title
	}
	if input.Description != nil {
		challenge.Description = *input.Description
	}
	if input.CategoryID != nil {
		challenge.CategoryID = *input.CategoryID
		category, err := s.categoryRepo.GetCategoryByID(challenge.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, exception.NewNotFoundException("Category", fmt.Sprintf("%d", challenge.CategoryID), "CATEGORY_NOT_FOUND")
		}
	}
	if input.MaxParticipants != nil {
		challenge.MaxParticipants = *input.MaxParticipants
	}
	if input.Visibility != nil {
		challenge.Visibility = *input.Visibility
	}
	if input.Rule != nil {
		challenge.Rule = *input.Rule
	}
	if input.ImageURL != nil {
		challenge.ImageURL = *input.ImageURL
	}
	if input.Timezone != nil {
		challenge.Timezone = *input.Timezone
	}
	if input.EndTime != nil {
		challenge.EndTime = input.EndTime
	}
	if input.CommentsEnabled != nil {
		challenge.CommentsEnabled = *input.CommentsEnabled
	}
	if input.IsStopped != nil {
		challenge.IsStopped = *input.IsStopped
	}

	updated, err := s.challengeRepo.UpdateChallenge(challenge)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *ChallengeService) DeleteChallenge(challengeID uint, currentUserID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, currentUserID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the creator can delete this challenge", "CHALLENGE_DELETE_FORBIDDEN")
	}

	return s.challengeRepo.DeleteChallenge(challengeID)
}

func (s *ChallengeService) StopChallenge(challengeID, currentUserID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, currentUserID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the challenge creator can stop this challenge", "STOP_CHALLENGE_NOT_AUTHORIZED")
	}

	return s.challengeRepo.StopChallenge(challengeID)
}

func (s *ChallengeService) GetChallengeByID(id, userID uint) (*dto.ChallengeDetailDTO, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(id, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}
	categoryName := ""
	if challenge.CategoryID > 0 {
		category, err := s.categoryRepo.GetCategoryByID(challenge.CategoryID)
		if err == nil && category != nil {
			categoryName = category.Name
		}
	}
	creatorUsername := ""
	if challenge.CreatorID > 0 {
		user, err := s.userRepo.GetUserByID(challenge.CreatorID)
		if err == nil && user != nil {
			creatorUsername = user.Username
		}
	}

	totalParticipants, err := s.participantRepo.GetParticipantCount(id)
	if err != nil {
		return nil, err
	}
	var mutualDTOs []dto.UserPreviewDTO
	if userID > 0 {
		mutualUsers, err := s.challengeRepo.GetMutualFollowersInChallenge(userID, id)
		if err == nil && len(mutualUsers) > 0 {
			mutualDTOs = make([]dto.UserPreviewDTO, len(mutualUsers))
			for i, u := range mutualUsers {
				mutualDTOs[i] = dto.UserPreviewDTO{
					ID:       u.ID,
					Username: u.Username,
				}
			}
		}
	}
	comments, err := s.commentRepo.GetChallengeComments(id, 0, 20)
	if err != nil {
		return nil, err
	}
	commentDTOs := make([]dto.CommentResponseDTO, len(comments))
	for i, c := range comments {
		username := ""
		if c.UserID > 0 {
			user, err := s.userRepo.GetUserByID(c.UserID)
			if err == nil && user != nil {
				username = user.Username
			}
		}
		commentDTOs[i] = dto.CommentResponseDTO{
			ID:       c.ID,
			UserID:   c.UserID,
			Username: username,
			Content:  c.Content,
		}
	}
	participants, err := s.participantRepo.GetParticipantsByChallenge(id, 0, 20)
	if err != nil {
		return nil, err
	}
	participantDTOs := make([]dto.ParticipantResponseDTO, len(participants))
	for i, p := range participants {
		username := ""
		if p.UserID > 0 {
			user, err := s.userRepo.GetUserByID(p.UserID)
			if err == nil && user != nil {
				username = user.Username
			}
		}
		participantDTOs[i] = dto.ParticipantResponseDTO{
			UserID:   p.UserID,
			Username: username,
		}
	}
	likeCount, _ := s.likeRepo.GetLikeCount(id)
	isUserLiked := false
	if userID > 0 {
		isUserLiked, _ = s.likeRepo.IsUserLikedChallenge(id, userID)
	}
	isUserParticipating := false
	if userID > 0 {
		isUserParticipating, err = s.participantRepo.IsUserParticipant(id, userID)
		if err != nil {
			return nil, err
		}
	}
	previewDTO := dto.ChallengePreviewDTO{
		ID:                  challenge.ID,
		Title:               challenge.Title,
		Description:         challenge.Description,
		Rule:                challenge.Rule,
		CategoryName:        categoryName,
		CreatorUsername:     creatorUsername,
		CreatorID:           challenge.CreatorID,
		Visibility:          challenge.Visibility,
		ImageURL:            challenge.ImageURL,
		MaxParticipants:     challenge.MaxParticipants,
		CurrentParticipants: int(totalParticipants),
		LikeCount:           likeCount,
		CommentCount:        uint(len(comments)),
		StartTime:           challenge.StartTime,
		EndTime:             challenge.EndTime,
		Timezone:            challenge.Timezone,
		CreatedAt:           challenge.CreatedAt,
		IsUserParticipating: isUserParticipating,
		IsUserLiked:         isUserLiked,
		MutualParticipants:  mutualDTOs,
	}
	return &dto.ChallengeDetailDTO{
		ChallengePreviewDTO: previewDTO,
		CommentsEnabled:     challenge.CommentsEnabled,
		Participants:        participantDTOs,
		Comments:            commentDTOs,
	}, nil
}

func (s *ChallengeService) GetAllChallenges(userID, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.GetAllChallenges(uint(userID), offset, limit)
}

func (s *ChallengeService) ListDiscoverableChallenges(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListDiscoverableChallenges(currentUserID, offset, limit)
}

func (s *ChallengeService) ListPublicChallenges(offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListPublicChallenges(offset, limit)
}

func (s *ChallengeService) ListChallengesByCreatorID(creatorID, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByCreatorID(creatorID, currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesByCreatorUsername(username string, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByCreatorUsername(username, currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesByCategoryID(categoryID, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByCategoryID(categoryID, currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesByCategoryName(name string, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByCategoryName(name, currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesByParticipantCount(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByParticipantCount(currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesByLikeCount(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesByLikeCount(currentUserID, offset, limit)
}

func (s *ChallengeService) ListChallengesStartingSoon(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesStartingSoon(currentUserID, offset, limit)
}

func (s *ChallengeService) ListTopCreatorsChallenge(offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListTopCreatorsChallenge(offset, limit)
}

func (s *ChallengeService) ListChallengesJoinedByUser(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesJoinedByUser(currentUserID, offset, limit)
}

func (s *ChallengeService) SearchChallenges(query string, visibility []enum.ChallengeVisibility, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.SearchChallenges(query, visibility, currentUserID, offset, limit)
}

func (s *ChallengeService) JoinPublicChallenge(userID, challengeID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.Visibility != enum.VisibilityPublic {
		return exception.NewForbiddenException("Challenge is not public", "CHALLENGE_NOT_PUBLIC")
	}

	if challenge.IsStopped {
		return exception.NewForbiddenException("Challenge is stopped", "CHALLENGE_STOPPED")
	}

	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isParticipant {
		return exception.NewConflictException("Participant", "user_id", "USER_ALREADY_PARTICIPANT")
	}
	participant := &model.ChallengeParticipant{
		ChallengeID: challengeID,
		UserID:      userID,
		Status:      enum.StatusJoined,
	}
	_, err = s.participantRepo.CreateParticipant(participant)
	return err
}

func (s *ChallengeService) JoinPrivateChallenge(userID, challengeID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.Visibility != enum.VisibilityPrivate {
		return exception.NewForbiddenException("Challenge is not private", "CHALLENGE_NOT_PRIVATE")
	}
	if challenge.IsStopped {
		return exception.NewForbiddenException("Challenge is stopped", "CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isParticipant {
		return exception.NewConflictException("Participant", "user_id", "USER_ALREADY_PARTICIPANT")
	}

	requests, err := s.joinRequestRepo.GetRequestsSentByUser(userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	for _, req := range requests {
		if req.ChallengeID == challengeID && req.Status == enum.RequestStatusPending {
			return exception.NewConflictException("JoinRequest", "challenge_id", "JOIN_REQUEST_ALREADY_EXISTS")
		}
	}

	// request := &model.ChallengeRequest{
	// 	ChallengeID: challengeID,
	// 	RequesterID: userID,
	// 	Status:      enum.StatusPending,
	// }
	_, err = s.joinRequestRepo.CreateRequest(userID, challengeID)
	return err
}

func (s *ChallengeService) InviteUserToChallenge(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, inviterID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != inviterID {
		return nil, exception.NewUnauthorizedException("Only the challenge creator can invite users", "INVITE_NOT_AUTHORIZED")
	}
	if challenge.IsStopped {
		return nil, exception.NewForbiddenException("Challenge is stopped", "CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return nil, err
	}
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if isParticipant {
		return nil, exception.NewConflictException("Participant", "user_id", "USER_ALREADY_PARTICIPANT")
	}

	invites, err := s.inviteRepo.GetInvitesSentToUser(inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	for _, invite := range invites {
		if invite.ChallengeID == challengeID && invite.Status == enum.RequestStatusPending {
			return nil, exception.NewConflictException("Invite", "challenge_id", "USER_ALREADY_INVITED")
		}
	}

	// invite := &model.ChallengeInvite{
	// 	ChallengeID: challengeID,
	// 	InviterID:   inviterID,
	// 	InviteeID:   inviteeID,
	// 	Status:      enum.StatusInvited,
	// }
	createdInvite, err := s.inviteRepo.CreateInvite(inviterID, challengeID, inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	return createdInvite, nil
}

func (s *ChallengeService) RemoveParticipant(challengeID, removerID, participantID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, removerID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != removerID {
		return exception.NewUnauthorizedException("Only the challenge creator can remove participants", "REMOVE_PARTICIPANT_NOT_AUTHORIZED")
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, participantID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return exception.NewBadRequestException("User is not a participant", "USER_NOT_PARTICIPANT", nil)
	}

	return s.participantRepo.DeleteParticipant(challengeID, participantID)
}

func (s *ChallengeService) ListChallengeParticipants(challengeID uint, userID uint, offset, limit int) ([]*model.ChallengeParticipant, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.Visibility == enum.VisibilityInvite || challenge.Visibility == enum.VisibilityPrivate {
		isCreator := challenge.CreatorID == userID
		isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
		if err != nil {
			return nil, err
		}
		if !isCreator && !isParticipant {
			return nil, exception.NewUnauthorizedException("You must be a participant or creator to view challenge participants", "PARTICIPANTS_ACCESS_DENIED")
		}
	}

	return s.participantRepo.GetParticipantsByChallenge(challengeID, offset, limit)
}

func (s *ChallengeService) AcceptJoinRequest(requestID, currentUserID uint) error {
	request, err := s.joinRequestRepo.GetRequest(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return exception.NewNotFoundException("Join request", fmt.Sprintf("%d", requestID), "JOIN_REQUEST_NOT_FOUND")
	}

	if request.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("Join request is not pending", "JOIN_REQUEST_NOT_PENDING", nil)
	}

	challenge, err := s.challengeRepo.GetChallengeByID(request.ChallengeID, currentUserID)
	if err != nil {
		return err
	}

	if challenge.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the challenge creator can accept join requests", "JOIN_REQUEST_NOT_AUTHORIZED")
	}

	if challenge.IsStopped {
		return exception.NewForbiddenException("Challenge is stopped", "CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(request.ChallengeID, challenge.MaxParticipants); err != nil {
		return err
	}
	request.Status = enum.RequestStatusAccepted
	_, err = s.joinRequestRepo.UpdateRequest(request)
	if err != nil {
		return err
	}

	participant := &model.ChallengeParticipant{
		ChallengeID: request.ChallengeID,
		UserID:      request.RequesterID,
		Status:      enum.StatusJoined,
	}
	_, err = s.participantRepo.CreateParticipant(participant)
	if err != nil {
		return err
	}
	err = s.joinRequestRepo.DeleteRequest(requestID)
	if err != nil {
		return err
	}
	return err
}

func (s *ChallengeService) DeclineJoinRequest(requestID, currentUserID uint) error {
	request, err := s.joinRequestRepo.GetRequest(requestID)
	if err != nil {
		return err
	}
	if request == nil {
		return exception.NewNotFoundException("Join request", fmt.Sprintf("%d", requestID), "JOIN_REQUEST_NOT_FOUND")
	}

	if request.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("Join request is not pending", "JOIN_REQUEST_NOT_PENDING", nil)
	}
	challenge, err := s.challengeRepo.GetChallengeByID(request.ChallengeID, currentUserID)
	if err != nil {
		return err
	}

	if challenge.CreatorID != currentUserID {
		return exception.NewUnauthorizedException("Only the challenge creator can decline join requests", "JOIN_REQUEST_NOT_AUTHORIZED")
	}
	request.Status = enum.RequestStatusRejected
	_, err = s.joinRequestRepo.UpdateRequest(request)
	if err != nil {
		return err
	}
	err = s.joinRequestRepo.DeleteRequest(requestID)
	if err != nil {
		return err
	}
	return err
}

func (s *ChallengeService) AcceptInvite(inviteID, currentUserID uint) error {
	invite, err := s.inviteRepo.GetInvite(inviteID)
	if err != nil {
		return err
	}
	if invite == nil {
		return exception.NewNotFoundException("Invite", fmt.Sprintf("%d", inviteID), "INVITE_NOT_FOUND")
	}
	if invite.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("Invite is not pending", "INVITE_NOT_PENDING", nil)
	}

	if invite.InviteeID != currentUserID {
		return exception.NewUnauthorizedException("You are not authorized to accept this invite", "INVITE_NOT_AUTHORIZED")
	}

	challenge, err := s.challengeRepo.GetChallengeByID(invite.ChallengeID, currentUserID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", invite.ChallengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.IsStopped {
		return exception.NewForbiddenException("Challenge is stopped", "CHALLENGE_STOPPED")
	}

	if err := s.checkParticipantLimit(invite.ChallengeID, challenge.MaxParticipants); err != nil {
		return err
	}

	invite.Status = enum.RequestStatusAccepted
	_, err = s.inviteRepo.UpdateInvite(invite)
	if err != nil {
		return err
	}
	participant := &model.ChallengeParticipant{
		ChallengeID: invite.ChallengeID,
		UserID:      currentUserID,
		Status:      enum.StatusJoined,
	}
	_, err = s.participantRepo.CreateParticipant(participant)
	if err != nil {
		return err
	}
	return s.inviteRepo.DeleteInvite(inviteID)
}

func (s *ChallengeService) DeclineInvite(inviteID, currentUserID uint) error {
	invite, err := s.inviteRepo.GetInvite(inviteID)
	if err != nil {
		return err
	}
	if invite == nil {
		return exception.NewNotFoundException("Invite", fmt.Sprintf("%d", inviteID), "INVITE_NOT_FOUND")
	}
	if invite.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("Invite is not pending", "INVITE_NOT_PENDING", nil)
	}
	if invite.InviteeID != currentUserID {
		return exception.NewUnauthorizedException("You are not authorized to decline this invite", "INVITE_NOT_AUTHORIZED")
	}
	invite.Status = enum.RequestStatusRejected
	_, err = s.inviteRepo.UpdateInvite(invite)
	if err != nil {
		return err
	}
	err = s.inviteRepo.DeleteInvite(inviteID)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChallengeService) LeaveChallenge(userID, challengeID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return err
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return exception.NewBadRequestException("User is not a participant", "USER_NOT_PARTICIPANT", nil)
	}

	return s.participantRepo.DeleteParticipant(challengeID, userID)
}

func (s *ChallengeService) AddComment(userID uint, input *dto.AddCommentDTO) (*model.ChallengeComment, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "COMMENT_ADD_BAD_INPUT", nil)
	}
	challenge, err := s.challengeRepo.GetChallengeByID(input.ChallengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.ChallengeID), "CHALLENGE_NOT_FOUND")
	}
	if !challenge.CommentsEnabled {
		return nil, exception.NewForbiddenException("Comments are disabled for this challenge", "COMMENTS_DISABLED")
	}
	isParticipant, err := s.participantRepo.IsUserParticipant(input.ChallengeID, userID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, exception.NewForbiddenException("Only participants can comment", "USER_NOT_PARTICIPANT")
	}

	comment := &model.ChallengeComment{
		ChallengeID: input.ChallengeID,
		UserID:      userID,
		Content:     input.Content,
	}
	createdComment, err := s.commentRepo.CreateComment(comment)
	return createdComment, err
}

func (s *ChallengeService) GetAllComments(challengeID, userID uint, offset, limit int) ([]*model.ChallengeComment, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if !challenge.CommentsEnabled {
		return nil, exception.NewForbiddenException("Comments are disabled for this challenge", "COMMENTS_DISABLED")
	}

	return s.commentRepo.GetChallengeComments(challengeID, offset, limit)
}

func (s *ChallengeService) GetComment(commentID, userID uint) (*dto.CommentResponseDTO, error) {
	comment, err := s.commentRepo.GetComment(commentID)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, exception.NewNotFoundException("Comment", fmt.Sprintf("%d", commentID), "COMMENT_NOT_FOUND")
	}
	challenge, err := s.challengeRepo.GetChallengeByID(comment.ChallengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewUnauthorizedException("You don't have permission to view this comment", "COMMENT_ACCESS_DENIED")
	}
	username := ""
	if comment.UserID > 0 {
		user, err := s.userRepo.GetUserByID(comment.UserID)
		if err == nil && user != nil {
			username = user.Username
		}
	}
	return &dto.CommentResponseDTO{
		ID:        comment.ID,
		UserID:    comment.UserID,
		Username:  username,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}, nil
}
func (s *ChallengeService) GetAllCategories() ([]*model.ChallengeCategoryModel, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *ChallengeService) GetRequestsSentByUser(userID uint) ([]*model.ChallengeRequest, error) {
	return s.joinRequestRepo.GetRequestsSentByUser(userID)
}

func (s *ChallengeService) GetInvitesSentToUser(userID uint) ([]*model.ChallengeInvite, error) {
	return s.inviteRepo.GetInvitesSentToUser(userID)
}

func (s *ChallengeService) GetInvitesSentFromChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeInvite, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, creatorID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != creatorID {
		return nil, exception.NewUnauthorizedException("Only the challenge creator can view invites", "INVITES_ACCESS_DENIED")
	}
	return s.inviteRepo.GetInvitesSentFromChallenge(challengeID)
}

func (s *ChallengeService) GetRequestsSentToChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeRequest, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, creatorID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != creatorID {
		return nil, exception.NewUnauthorizedException("Only the challenge creator can view join requests", "REQUESTS_ACCESS_DENIED")
	}

	return s.joinRequestRepo.GetRequestsSentToChallenge(challengeID)
}

func (s *ChallengeService) IsUserParticipant(challengeID, userID uint) (bool, error) {
	return s.participantRepo.IsUserParticipant(challengeID, userID)
}

func (s *ChallengeService) IsChallengeCreator(challengeID, userID uint) (bool, error) {
	return s.challengeRepo.IsChallengeCreator(challengeID, userID)
}

func (s *ChallengeService) GetChallengeParticipantCount(challengeID uint) (int, error) {
	return s.participantRepo.GetParticipantCount(challengeID)
}

func (s *ChallengeService) GetChallengesUserIsParticipating(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	return s.challengeRepo.ListChallengesByParticipant(userID, offset, limit)
}
func (s *ChallengeService) GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	return s.challengeRepo.GetMutualFollowersInChallenge(userID, challengeID)
}

func (s *ChallengeService) LikeChallenge(userID, challengeID uint) error {
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return exception.NewForbiddenException("Only participants can like challenges", "USER_NOT_PARTICIPANT")
	}

	isLiked, err := s.likeRepo.IsUserLikedChallenge(challengeID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return exception.NewConflictException("Like", "user_id", "USER_ALREADY_LIKED")
	}

	like := &model.ChallengeLike{
		ChallengeID: challengeID,
		UserID:      userID,
	}

	return s.likeRepo.CreateLike(like)
}

func (s *ChallengeService) UnlikeChallenge(userID, challengeID uint) error {
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return err
	}
	if !isParticipant {
		return exception.NewForbiddenException("Only participants can unlike challenges", "USER_NOT_PARTICIPANT")
	}

	isLiked, err := s.likeRepo.IsUserLikedChallenge(challengeID, userID)
	if err != nil {
		return err
	}
	if !isLiked {
		return exception.NewBadRequestException("User has not liked this challenge", "USER_NOT_LIKED", nil)
	}

	return s.likeRepo.DeleteLike(challengeID, userID)
}

func (s *ChallengeService) IsUserLikedChallenge(userID, challengeID uint) (bool, error) {
	return s.likeRepo.IsUserLikedChallenge(challengeID, userID)
}

func (s *ChallengeService) GetChallengeLikeCount(challengeID uint) (uint, error) {
	return s.likeRepo.GetLikeCount(challengeID)
}

// Helper
func (s *ChallengeService) checkParticipantLimit(challengeID uint, maxParticipants uint) error {
	if maxParticipants > 0 {
		count, err := s.participantRepo.GetParticipantCount(challengeID)
		if err != nil {
			return err
		}
		if count >= int(maxParticipants) {
			return exception.NewForbiddenException("Challenge participant limit reached", "CHALLENGE_PARTICIPANT_LIMIT_REACHED")
		}
	}
	return nil
}
