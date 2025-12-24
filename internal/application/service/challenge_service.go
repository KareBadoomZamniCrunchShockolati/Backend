package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/storage"
)

type ChallengeService struct {
	challengeRepo   repository.ChallengeRepository
	participantRepo repository.ChallengeParticipantRepository
	commentRepo     repository.CommentRepository
	inviteRepo      repository.ChallengeInviteRepository
	joinRequestRepo repository.ChallengeJoinRequestRepository
	categoryRepo    repository.CategoryRepository
	userRepo        repository.UserRepository
	followRepo      repository.FollowRepository
	likeRepo        repository.LikeRepository
	objectStorage   storage.ObjectStorage
}

func NewChallengeService(
	challengeRepo repository.ChallengeRepository,
	participantRepo repository.ChallengeParticipantRepository,
	commentRepo repository.CommentRepository,
	inviteRepo repository.ChallengeInviteRepository,
	joinRequestRepo repository.ChallengeJoinRequestRepository,
	categoryRepo repository.CategoryRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	likeRepo repository.LikeRepository,
	objectStorage storage.ObjectStorage,
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
		objectStorage:   objectStorage,
	}
}

// Validation helpers for map coordinates
func isValidLatitude(lat float64) bool { return lat >= -90 && lat <= 90 }
func isValidLongitude(lng float64) bool { return lng >= -180 && lng <= 180 }

// Challenge CRUD operations
func (s *ChallengeService) CreateChallenge(input *dto.CreateChallengeDTO) (*model.ChallengeModel, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("CHALLENGE_CREATE_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}

	// Validate map coordinates if provided (must be provided together)
	if input.Latitude != nil && input.Longitude != nil {
		if !isValidLatitude(*input.Latitude) || !isValidLongitude(*input.Longitude) {
			return nil, exception.NewBadRequestException("INVALID_MAP_COORDINATES", map[string]any{
				"reason":    "invalid_range",
				"latitude":  *input.Latitude,
				"longitude": *input.Longitude,
				"lat_min":   -90,
				"lat_max":   90,
				"lng_min":   -180,
				"lng_max":   180,
			})
		}
	} else if input.Latitude != nil || input.Longitude != nil {
		return nil, exception.NewBadRequestException("MISSING_MAP_COORDINATES", map[string]any{
			"reason":            "one_coordinate_missing",
			"latitude_present":  input.Latitude != nil,
			"longitude_present": input.Longitude != nil,
		})
	}

	category, err := s.categoryRepo.GetCategoryByID(input.CategoryID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if category == nil {
		return nil, exception.NewNotFoundException("Category", fmt.Sprintf("%d", input.CategoryID), "CATEGORY_NOT_FOUND")
	}
	user, err := s.userRepo.GetUserByID(input.CreatorID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", input.CreatorID), "USER_NOT_FOUND")
	}

	goal := 1
	if input.Goal != nil {
		if *input.Goal < 1 {
			return nil, exception.NewBadRequestException("INVALID_GOAL", map[string]any{
				"reason": "goal_must_be_at_least_1",
				"min":    1,
				"value":  *input.Goal,
			})
		}
		goal = *input.Goal
	}

	challenge := &model.ChallengeModel{
		Title:           input.Title,
		Description:     input.Description,
		CategoryID:      input.CategoryID,
		CreatorID:       input.CreatorID,
		MaxParticipants: input.MaxParticipants,
		Visibility:      input.Visibility,

		Latitude:  0,
		Longitude: 0,
		Address:   "",

		Goal:            goal,
		Rule:            input.Rule,
		StartTime:        input.StartTime,
		EndTime:          &input.EndTime,
		Timezone:         input.Timezone,
		IsStopped:        false,
		CommentsEnabled:  input.CommentsEnabled,
	}

	// Set location if provided
	if input.Latitude != nil && input.Longitude != nil {
		challenge.Latitude = *input.Latitude
		challenge.Longitude = *input.Longitude
	}
	if input.Address != nil {
		challenge.Address = *input.Address
	}

	created, err := s.challengeRepo.CreateChallenge(challenge)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return created, nil
}

func (s *ChallengeService) UpdateChallenge(id uint, currentUserID uint, input *dto.UpdateChallengeDTO) (*model.ChallengeModel, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(id, currentUserID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", id), "CHALLENGE_NOT_FOUND")
	}

	if challenge.CreatorID != currentUserID {
		return nil, exception.NewForbiddenException("CHALLENGE_UPDATE_FORBIDDEN")
	}
	if input.StartTime != nil {
		if challenge.StartTime.Before(time.Now()) {
			return nil, exception.NewBadRequestException("CHALLENGE_ALREADY_STARTED", map[string]any{
				"reason": "cannot_change_start_time_after_started",
			})
		}
		challenge.StartTime = *input.StartTime
	}
	if !challenge.StartTime.Before(time.Now()) {
		if input.Goal != nil {
			if *input.Goal < 1 {
				return nil, exception.NewBadRequestException("INVALID_GOAL", map[string]any{
					"reason": "goal_must_be_at_least_1",
					"min":    1,
					"value":  *input.Goal,
				})
			}
			challenge.Goal = *input.Goal
		}
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
			return nil, exception.NewRepositoryError(err)
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

	if input.Latitude != nil || input.Longitude != nil {
		if (input.Latitude != nil && input.Longitude == nil) || (input.Latitude == nil && input.Longitude != nil) {
			return nil, exception.NewBadRequestException("MISSING_MAP_COORDINATES", map[string]any{
				"reason":            "one_coordinate_missing",
				"latitude_present":  input.Latitude != nil,
				"longitude_present": input.Longitude != nil,
			})
		}

		if input.Latitude != nil && input.Longitude != nil {
			if !isValidLatitude(*input.Latitude) || !isValidLongitude(*input.Longitude) {
				return nil, exception.NewBadRequestException("INVALID_MAP_COORDINATES", map[string]any{
					"reason":    "invalid_range",
					"latitude":  *input.Latitude,
					"longitude": *input.Longitude,
					"lat_min":   -90,
					"lat_max":   90,
					"lng_min":   -180,
					"lng_max":   180,
				})
			}
			challenge.Latitude = *input.Latitude
			challenge.Longitude = *input.Longitude
		}
	}

	if input.Address != nil {
		challenge.Address = *input.Address
	}
	if input.Rule != nil {
		challenge.Rule = *input.Rule
	}
	if input.CoverImage != nil {
		challenge.CoverImage = *input.CoverImage
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
		return nil, exception.NewRepositoryError(err)
	}
	return updated, nil
}

func (s *ChallengeService) DeleteChallenge(challengeID uint, currentUserID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, currentUserID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != currentUserID {
		return exception.NewForbiddenException("CHALLENGE_DELETE_FORBIDDEN")
	}
	return s.challengeRepo.DeleteChallenge(challengeID)
}

func (s *ChallengeService) StopChallenge(challengeID, currentUserID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, currentUserID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != currentUserID {
		return exception.NewForbiddenException("STOP_CHALLENGE_NOT_AUTHORIZED")
	}
	return s.challengeRepo.StopChallenge(challengeID)
}

func (s *ChallengeService) buildNestedComments(comments []*model.Comment, userID uint) []*dto.CommentResponseDTO {
	commentMap := make(map[uint]*dto.CommentResponseDTO)
	var rootComments []*dto.CommentResponseDTO

	for _, comment := range comments {
		username := ""
		if comment.UserID > 0 {
			user, err := s.userRepo.GetUserByID(comment.UserID)
			if err == nil && user != nil {
				username = user.Username
			}
		}

		likeCount, _ := s.likeRepo.GetLikeCount(model.LikeTypeComment, comment.ID)
		isLiked, _ := s.likeRepo.IsUserLiked(model.LikeTypeComment, comment.ID, userID)

		commentDTO := &dto.CommentResponseDTO{
			ID:         comment.ID,
			EntityType: string(comment.EntityType),
			EntityID:   comment.EntityID,
			UserID:     comment.UserID,
			Username:   username,
			Content:    comment.Content,
			ParentID:   comment.ParentID,
			LikeCount:  likeCount,
			IsLiked:    isLiked,
			CreatedAt:  comment.CreatedAt,
			Replies:    []*dto.CommentResponseDTO{},
		}
		commentMap[comment.ID] = commentDTO

		if comment.ParentID == nil || *comment.ParentID == 0 {
			rootComments = append(rootComments, commentDTO)
		}
	}

	for _, comment := range comments {
		if comment.ParentID != nil && *comment.ParentID != 0 {
			if parent, exists := commentMap[*comment.ParentID]; exists {
				parent.Replies = append(parent.Replies, commentMap[comment.ID])
			}
		}
	}

	return rootComments
}

func (s *ChallengeService) GetChallengeByID(id, userID uint) (*dto.ChallengeDetailDTO, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(id, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
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
		return nil, exception.NewRepositoryError(err)
	}
	var mutualDTOs []dto.UserPreviewDTO
	if userID > 0 {
		mutualUsers, err := s.challengeRepo.GetMutualFollowersInChallenge(userID, id)
		if err == nil && len(mutualUsers) > 0 {
			mutualDTOs = make([]dto.UserPreviewDTO, len(mutualUsers))
			for i, u := range mutualUsers {
				mutualDTOs[i] = dto.UserPreviewDTO{ID: u.ID, Username: u.Username}
			}
		}
	}

	comments, err := s.commentRepo.GetComments(model.CommentTypeChallenge, id, 0, 20)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	commentDTOs := s.buildNestedComments(comments, userID)

	participants, err := s.participantRepo.GetParticipantsByChallenge(id, 0, 20)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
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
		participantDTOs[i] = dto.ParticipantResponseDTO{UserID: p.UserID, Username: username}
	}

	likeCount, _ := s.likeRepo.GetLikeCount(model.LikeTypeChallenge, id)
	isUserLiked := false
	if userID > 0 {
		isUserLiked, _ = s.likeRepo.IsUserLiked(model.LikeTypeChallenge, id, userID)
	}

	isUserParticipating := false
	if userID > 0 {
		isUserParticipating, err = s.participantRepo.IsUserParticipant(id, userID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
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
		Latitude:            challenge.Latitude,
		Longitude:           challenge.Longitude,
		Address:             challenge.Address,
		Goal:                challenge.Goal,
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
		CoverImage:          challenge.CoverImage,
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
func (s *ChallengeService) ListTopCreators(offset, limit int) ([]*dto.TopCreatorDTO, error) {
	return s.challengeRepo.ListTopCreators(offset, limit)
}
func (s *ChallengeService) ListChallengesJoinedByUser(currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.ListChallengesJoinedByUser(currentUserID, offset, limit)
}
func (s *ChallengeService) SearchChallenges(query string, visibility []enum.ChallengeVisibility, currentUserID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.SearchChallenges(query, visibility, currentUserID, offset, limit)
}
func (s *ChallengeService) SearchChallengesUserIsParticipating(userID uint, query string, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return s.challengeRepo.SearchChallengesUserIsParticipating(userID, query, offset, limit)
}

func (s *ChallengeService) JoinPublicChallenge(userID, challengeID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.Visibility != enum.VisibilityPublic {
		return exception.NewForbiddenException("CHALLENGE_NOT_PUBLIC")
	}
	if challenge.IsStopped {
		return exception.NewForbiddenException("CHALLENGE_STOPPED")
	}

	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isParticipant {
		return exception.NewConflictException("USER_ALREADY_PARTICIPANT", "Participant", "user_id", fmt.Sprintf("%d", userID))
	}
	participant := &model.ChallengeParticipant{
		ChallengeID: challengeID,
		UserID:      userID,
		Status:      enum.StatusJoined,
	}
	_, err = s.participantRepo.CreateParticipant(participant)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
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
		return exception.NewForbiddenException("CHALLENGE_NOT_PRIVATE")
	}
	if challenge.IsStopped {
		return exception.NewForbiddenException("CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isParticipant {
		return exception.NewConflictException("USER_ALREADY_PARTICIPANT", "Participant", "user_id", fmt.Sprintf("%d", userID))
	}

	requests, err := s.joinRequestRepo.GetRequestsSentByUser(userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	for _, req := range requests {
		if req.ChallengeID == challengeID && req.Status == enum.RequestStatusPending {
			return exception.NewConflictException("JOIN_REQUEST_ALREADY_EXISTS", "JoinRequest", "challenge_id", fmt.Sprintf("%d", challengeID))
		}
	}

	_, err = s.joinRequestRepo.CreateRequest(userID, challengeID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) InviteUserToChallenge(inviterID, challengeID, inviteeID uint) (*model.ChallengeInvite, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, inviterID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != inviterID {
		return nil, exception.NewForbiddenException("INVITE_NOT_AUTHORIZED")
	}
	if challenge.IsStopped {
		return nil, exception.NewForbiddenException("CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(challengeID, challenge.MaxParticipants); err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if isParticipant {
		return nil, exception.NewConflictException("USER_ALREADY_PARTICIPANT", "Participant", "user_id", fmt.Sprintf("%d", inviteeID))
	}

	invites, err := s.inviteRepo.GetInvitesSentToUser(inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	for _, inv := range invites {
		if inv.ChallengeID == challengeID && inv.Status == enum.RequestStatusPending {
			return nil, exception.NewConflictException("USER_ALREADY_INVITED", "Invite", "challenge_id", fmt.Sprintf("%d", challengeID))
		}
	}

	createdInvite, err := s.inviteRepo.CreateInvite(inviterID, challengeID, inviteeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return createdInvite, nil
}

func (s *ChallengeService) RemoveParticipant(challengeID, removerID, participantID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, removerID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != removerID {
		return exception.NewForbiddenException("REMOVE_PARTICIPANT_NOT_AUTHORIZED")
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, participantID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isParticipant {
		return exception.NewBadRequestException("USER_NOT_PARTICIPANT", map[string]any{
			"reason":  "target_user_is_not_participant",
			"user_id": participantID,
		})
	}

	err = s.participantRepo.DeleteParticipant(challengeID, participantID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) ListChallengeParticipants(challengeID uint, userID uint, offset, limit int) ([]*model.ChallengeParticipant, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	if challenge.Visibility == enum.VisibilityInvite || challenge.Visibility == enum.VisibilityPrivate {
		isCreator := challenge.CreatorID == userID
		isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		if !isCreator && !isParticipant {
			return nil, exception.NewForbiddenException("PARTICIPANTS_ACCESS_DENIED")
		}
	}

	participants, err := s.participantRepo.GetParticipantsByChallenge(challengeID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return participants, nil
}

func (s *ChallengeService) AcceptJoinRequest(requestID, currentUserID uint) error {
	request, err := s.joinRequestRepo.GetRequest(requestID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if request == nil {
		return exception.NewNotFoundException("Join request", fmt.Sprintf("%d", requestID), "JOIN_REQUEST_NOT_FOUND")
	}
	if request.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("JOIN_REQUEST_NOT_PENDING", map[string]any{
			"request_id": requestID,
			"status":     request.Status,
		})
	}

	challenge, err := s.challengeRepo.GetChallengeByID(request.ChallengeID, currentUserID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge.CreatorID != currentUserID {
		return exception.NewForbiddenException("JOIN_REQUEST_NOT_AUTHORIZED")
	}
	if challenge.IsStopped {
		return exception.NewForbiddenException("CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(request.ChallengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}
	request.Status = enum.RequestStatusAccepted
	if _, err := s.joinRequestRepo.UpdateRequest(request); err != nil {
		return exception.NewRepositoryError(err)
	}

	participant := &model.ChallengeParticipant{
		ChallengeID: request.ChallengeID,
		UserID:      request.RequesterID,
		Status:      enum.StatusJoined,
	}
	if _, err := s.participantRepo.CreateParticipant(participant); err != nil {
		return exception.NewRepositoryError(err)
	}
	err = s.joinRequestRepo.DeleteRequest(requestID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) DeclineJoinRequest(requestID, currentUserID uint) error {
	request, err := s.joinRequestRepo.GetRequest(requestID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if request == nil {
		return exception.NewNotFoundException("Join request", fmt.Sprintf("%d", requestID), "JOIN_REQUEST_NOT_FOUND")
	}
	if request.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("JOIN_REQUEST_NOT_PENDING", map[string]any{
			"request_id": requestID,
			"status":     request.Status,
		})
	}
	challenge, err := s.challengeRepo.GetChallengeByID(request.ChallengeID, currentUserID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge.CreatorID != currentUserID {
		return exception.NewForbiddenException("JOIN_REQUEST_NOT_AUTHORIZED")
	}
	request.Status = enum.RequestStatusRejected
	if _, err := s.joinRequestRepo.UpdateRequest(request); err != nil {
		return exception.NewRepositoryError(err)
	}
	err = s.joinRequestRepo.DeleteRequest(requestID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) AcceptInvite(inviteID, currentUserID uint) error {
	invite, err := s.inviteRepo.GetInvite(inviteID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if invite == nil {
		return exception.NewNotFoundException("Invite", fmt.Sprintf("%d", inviteID), "INVITE_NOT_FOUND")
	}
	if invite.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("INVITE_NOT_PENDING", map[string]any{
			"invite_id": inviteID,
			"status":    invite.Status,
		})
	}
	if invite.InviteeID != currentUserID {
		return exception.NewForbiddenException("INVITE_NOT_AUTHORIZED")
	}

	challenge, err := s.challengeRepo.GetChallengeByID(invite.ChallengeID, currentUserID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", invite.ChallengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.IsStopped {
		return exception.NewForbiddenException("CHALLENGE_STOPPED")
	}
	if err := s.checkParticipantLimit(invite.ChallengeID, challenge.MaxParticipants); err != nil {
		return exception.NewRepositoryError(err)
	}

	invite.Status = enum.RequestStatusAccepted
	if _, err := s.inviteRepo.UpdateInvite(invite); err != nil {
		return exception.NewRepositoryError(err)
	}
	participant := &model.ChallengeParticipant{
		ChallengeID: invite.ChallengeID,
		UserID:      currentUserID,
		Status:      enum.StatusJoined,
	}
	if _, err := s.participantRepo.CreateParticipant(participant); err != nil {
		return exception.NewRepositoryError(err)
	}

	err = s.inviteRepo.DeleteInvite(inviteID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) DeclineInvite(inviteID, currentUserID uint) error {
	invite, err := s.inviteRepo.GetInvite(inviteID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if invite == nil {
		return exception.NewNotFoundException("Invite", fmt.Sprintf("%d", inviteID), "INVITE_NOT_FOUND")
	}
	if invite.Status != enum.RequestStatusPending {
		return exception.NewBadRequestException("INVITE_NOT_PENDING", map[string]any{
			"invite_id": inviteID,
			"status":    invite.Status,
		})
	}
	if invite.InviteeID != currentUserID {
		return exception.NewForbiddenException("INVITE_NOT_AUTHORIZED")
	}
	invite.Status = enum.RequestStatusRejected
	if _, err := s.inviteRepo.UpdateInvite(invite); err != nil {
		return exception.NewRepositoryError(err)
	}
	err = s.inviteRepo.DeleteInvite(inviteID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) LeaveChallenge(userID, challengeID uint) error {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isParticipant {
		return exception.NewBadRequestException("USER_NOT_PARTICIPANT", map[string]any{
			"user_id":      userID,
			"challenge_id": challengeID,
		})
	}

	err = s.participantRepo.DeleteParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) AddComment(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error) {
	return s.AddCommentToChallenge(userID, input)
}

func (s *ChallengeService) AddCommentToChallenge(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("COMMENT_ADD_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}
	if input.EntityType != "challenge" {
		return nil, exception.NewBadRequestException("INVALID_ENTITY_TYPE", map[string]any{
			"expected": "challenge",
			"got":      input.EntityType,
		})
	}

	challenge, err := s.challengeRepo.GetChallengeByID(input.EntityID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.EntityID), "CHALLENGE_NOT_FOUND")
	}
	if !challenge.CommentsEnabled {
		return nil, exception.NewForbiddenException("COMMENTS_DISABLED")
	}
	isParticipant, err := s.participantRepo.IsUserParticipant(input.EntityID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if !isParticipant {
		return nil, exception.NewForbiddenException("USER_NOT_PARTICIPANT")
	}

	comment := &model.Comment{
		EntityType: model.CommentTypeChallenge,
		EntityID:   input.EntityID,
		UserID:     userID,
		Content:    input.Content,
		ParentID:   input.ParentID,
	}
	createdComment, err := s.commentRepo.CreateComment(comment)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return createdComment, nil
}

func (s *ChallengeService) GetAllComments(challengeID, userID uint, offset, limit int) ([]*dto.CommentResponseDTO, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if !challenge.CommentsEnabled {
		return nil, exception.NewForbiddenException("COMMENTS_DISABLED")
	}

	comments, err := s.commentRepo.GetComments(model.CommentTypeChallenge, challengeID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	return s.buildNestedComments(comments, userID), nil
}

func (s *ChallengeService) GetComment(commentID, userID uint) (*dto.CommentResponseDTO, error) {
	comment, err := s.commentRepo.GetComment(commentID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if comment == nil {
		return nil, exception.NewNotFoundException("Comment", fmt.Sprintf("%d", commentID), "COMMENT_NOT_FOUND")
	}

	if comment.EntityType != model.CommentTypeChallenge {
		return nil, exception.NewForbiddenException("COMMENT_ACCESS_DENIED")
	}

	challenge, err := s.challengeRepo.GetChallengeByID(comment.EntityID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewForbiddenException("COMMENT_ACCESS_DENIED")
	}
	username := ""
	if comment.UserID > 0 {
		user, err := s.userRepo.GetUserByID(comment.UserID)
		if err == nil && user != nil {
			username = user.Username
		}
	}

	likeCount, _ := s.likeRepo.GetLikeCount(model.LikeTypeComment, comment.ID)
	isLiked, _ := s.likeRepo.IsUserLiked(model.LikeTypeComment, comment.ID, userID)

	return &dto.CommentResponseDTO{
		ID:         comment.ID,
		EntityType: string(comment.EntityType),
		EntityID:   comment.EntityID,
		UserID:     comment.UserID,
		Username:   username,
		Content:    comment.Content,
		ParentID:   comment.ParentID,
		LikeCount:  likeCount,
		IsLiked:    isLiked,
		CreatedAt:  comment.CreatedAt,
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
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != creatorID {
		return nil, exception.NewForbiddenException("INVITES_ACCESS_DENIED")
	}
	invites, err := s.inviteRepo.GetInvitesSentFromChallenge(challengeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return invites, nil
}

func (s *ChallengeService) GetRequestsSentToChallenge(challengeID uint, creatorID uint) ([]*model.ChallengeRequest, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, creatorID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if challenge.CreatorID != creatorID {
		return nil, exception.NewForbiddenException("REQUESTS_ACCESS_DENIED")
	}
	requests, err := s.joinRequestRepo.GetRequestsSentToChallenge(challengeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return requests, nil
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
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	users, err := s.challengeRepo.GetMutualFollowersInChallenge(userID, challengeID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return users, nil
}

func (s *ChallengeService) LikeChallenge(userID, challengeID uint) error {
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isParticipant {
		return exception.NewForbiddenException("USER_NOT_PARTICIPANT")
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeTypeChallenge, challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isLiked {
		return exception.NewConflictException("USER_ALREADY_LIKED", "Like", "user_id", fmt.Sprintf("%d", userID))
	}

	like := &model.Like{
		EntityType: model.LikeTypeChallenge,
		EntityID:   challengeID,
		UserID:     userID,
	}
	err = s.likeRepo.CreateLike(like)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) UnlikeChallenge(userID, challengeID uint) error {
	isParticipant, err := s.participantRepo.IsUserParticipant(challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isParticipant {
		return exception.NewForbiddenException("USER_NOT_PARTICIPANT")
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeTypeChallenge, challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isLiked {
		return exception.NewBadRequestException("USER_NOT_LIKED", map[string]any{
			"user_id":      userID,
			"challenge_id": challengeID,
		})
	}

	err = s.likeRepo.DeleteLike(model.LikeTypeChallenge, challengeID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *ChallengeService) IsUserLikedChallenge(userID, challengeID uint) (bool, error) {
	return s.likeRepo.IsUserLiked(model.LikeTypeChallenge, challengeID, userID)
}

func (s *ChallengeService) GetChallengeLikeCount(challengeID uint) (uint, error) {
	return s.likeRepo.GetLikeCount(model.LikeTypeChallenge, challengeID)
}

// Helper
func (s *ChallengeService) checkParticipantLimit(challengeID uint, maxParticipants uint) error {
	if maxParticipants > 0 {
		count, err := s.participantRepo.GetParticipantCount(challengeID)
		if err != nil {
			return err
		}
		if count >= int(maxParticipants) {
			return exception.NewForbiddenException("CHALLENGE_PARTICIPANT_LIMIT_REACHED")
		}
	}
	return nil
}

func (s *ChallengeService) UploadChallengeCover(ctx context.Context, userID uint, challengeID uint, data []byte, mime string) (*model.ChallengeModel, error) {
	ch, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if ch == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}
	if ch.CreatorID != userID {
		return nil, exception.NewForbiddenException("COVER_FORBIDDEN")
	}

	ext := ".jpg"
	if mime == "image/png" {
		ext = ".png"
	} else if mime == "image/webp" {
		ext = ".webp"
	}
	key := fmt.Sprintf("challenges/%d/covers/%d%s", challengeID, time.Now().UTC().UnixNano(), ext)
	publicURL, err := s.objectStorage.Upload(ctx, key, mime, bytes.NewReader(data))
	if err != nil {
		return nil, exception.NewInternalServerException("S3_UPLOAD_FAILED", map[string]any{
			"challenge_id": challengeID,
			"key":          key,
			"mime":         mime,
		}, err)
	}
	ch.CoverImage = publicURL
	updated, err := s.challengeRepo.UpdateChallenge(ch)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return updated, nil
}
