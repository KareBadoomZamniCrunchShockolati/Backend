package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"fmt"
)

type PostService struct {
	postRepo      repository.PostRepository
	commentRepo   repository.CommentRepository
	likeRepo      repository.LikeRepository
	userRepo      repository.UserRepository
	challengeRepo repository.ChallengeRepository
}

func NewPostService(
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	userRepo repository.UserRepository,
	challengeRepo repository.ChallengeRepository,
) *PostService {
	return &PostService{
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		likeRepo:      likeRepo,
		userRepo:      userRepo,
		challengeRepo: challengeRepo,
	}
}

func (s *PostService) CreatePost(userID uint, input *dto.CreatePostDTO) (*model.Post, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "POST_CREATE_BAD_INPUT", nil)
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND")
	}

	if input.ChallengeID != nil {
		challenge, err := s.challengeRepo.GetChallengeByID(*input.ChallengeID, userID)
		if err != nil {
			return nil, err
		}
		if challenge == nil {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", *input.ChallengeID), "CHALLENGE_NOT_FOUND")
		}
	}

	if len(input.Pictures) > 5 {
		return nil, exception.NewBadRequestException("Maximum 5 pictures allowed", "POST_TOO_MANY_PICTURES", nil)
	}

	post := &model.Post{
		UserID:      userID,
		Description: input.Description,
		ChallengeID: input.ChallengeID,
		Pictures:    input.Pictures,
	}

	return s.postRepo.CreatePost(post)
}

func (s *PostService) GetPost(postID uint, userID uint) (*dto.PostResponseDTO, error) {
	// Get the post
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	// Get user info
	user, err := s.userRepo.GetUserByID(post.UserID)
	if err != nil {
		return nil, err
	}
	username := ""
	if user != nil {
		username = user.Username
	}

	// Get like count
	likeCount, err := s.likeRepo.GetLikeCount(model.LikeTypePost, post.ID)
	if err != nil {
		likeCount = 0
	}

	// Get comment count
	commentCount, err := s.commentRepo.GetCommentCount(model.CommentTypePost, post.ID)
	if err != nil {
		commentCount = 0
	}

	// Check if current user liked this post
	isLiked, err := s.likeRepo.IsUserLiked(model.LikeTypePost, post.ID, userID)
	if err != nil {
		isLiked = false
	}

	// Build response
	response := &dto.PostResponseDTO{
		ID:           post.ID,
		UserID:       post.UserID,
		Username:     username,
		Description:  post.Description,
		ChallengeID:  post.ChallengeID,
		Pictures:     post.Pictures,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		IsLiked:      isLiked,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}

	return response, nil
}

func (s *PostService) UpdatePost(postID, userID uint, input *dto.UpdatePostDTO) (*model.Post, error) {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	if post.UserID != userID {
		return nil, exception.NewUnauthorizedException("Only the post creator can update this post", "POST_UPDATE_FORBIDDEN")
	}

	if input.Description != nil {
		post.Description = *input.Description
	}
	if input.Pictures != nil {
		if len(*input.Pictures) > 5 {
			return nil, exception.NewBadRequestException("Maximum 5 pictures allowed", "POST_TOO_MANY_PICTURES", nil)
		}
		post.Pictures = *input.Pictures
	}

	return s.postRepo.UpdatePost(post)
}

func (s *PostService) DeletePost(postID, userID uint) error {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	if post.UserID != userID {
		return exception.NewUnauthorizedException("Only the post creator can delete this post", "POST_DELETE_FORBIDDEN")
	}

	return s.postRepo.DeletePost(postID)
}

func (s *PostService) GetUserPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	posts, err := s.postRepo.GetPostsByUser(userID, offset, limit)
	if err != nil {
		return nil, err
	}
	return s.enrichPostsWithDetails(posts, userID)
}

func (s *PostService) GetFeedPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	posts, err := s.postRepo.GetFeedPosts(userID, offset, limit)
	if err != nil {
		return nil, err
	}
	return s.enrichPostsWithDetails(posts, userID)
}

func (s *PostService) AddComment(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("Input cannot be nil", "COMMENT_ADD_BAD_INPUT", nil)
	}

	switch input.EntityType {
	case "challenge":
		challenge, err := s.challengeRepo.GetChallengeByID(input.EntityID, userID)
		if err != nil {
			return nil, err
		}
		if challenge == nil {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.EntityID), "CHALLENGE_NOT_FOUND")
		}
	case "post":
		post, err := s.postRepo.GetPost(input.EntityID)
		if err != nil {
			return nil, err
		}
		if post == nil {
			return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", input.EntityID), "POST_NOT_FOUND")
		}
	default:
		return nil, exception.NewBadRequestException("Invalid entity type", "INVALID_ENTITY_TYPE", nil)
	}

	comment := &model.Comment{
		EntityType: model.CommentType(input.EntityType),
		EntityID:   input.EntityID,
		UserID:     userID,
		Content:    input.Content,
		ParentID:   input.ParentID,
	}

	return s.commentRepo.CreateComment(comment)
}

func (s *PostService) GetComments(entityType string, entityID, userID uint, offset, limit int) ([]*dto.CommentResponseDTO, error) {
	comments, err := s.commentRepo.GetComments(model.CommentType(entityType), entityID, offset, limit)
	if err != nil {
		return nil, err
	}

	nestedComments := s.buildCommentTree(comments)
	return s.convertCommentsToDTO(nestedComments, userID), nil
}

func (s *PostService) buildCommentTree(comments []*model.Comment) []*model.Comment {
	commentMap := make(map[uint]*model.Comment)
	var rootComments []*model.Comment

	// First pass: create map and identify root comments
	for _, comment := range comments {
		commentCopy := *comment
		commentMap[commentCopy.ID] = &commentCopy

		if commentCopy.ParentID == nil || *commentCopy.ParentID == 0 {
			rootComments = append(rootComments, &commentCopy)
		}
	}

	// Second pass: build the tree structure
	for _, comment := range comments {
		if comment.ParentID != nil && *comment.ParentID != 0 {
			if parent, exists := commentMap[*comment.ParentID]; exists {
				parent.AddReply(commentMap[comment.ID])
			}
		}
	}

	return rootComments
}

// convertCommentsToDTO converts model comments to DTO with nested structure
func (s *PostService) convertCommentsToDTO(comments []*model.Comment, userID uint) []*dto.CommentResponseDTO {
	if comments == nil {
		return []*dto.CommentResponseDTO{}
	}

	dtos := make([]*dto.CommentResponseDTO, len(comments))
	for i, comment := range comments {
		dtos[i] = s.convertCommentToDTO(comment, userID)
	}
	return dtos
}

// convertCommentToDTO converts a single comment to DTO with nested replies
func (s *PostService) convertCommentToDTO(comment *model.Comment, userID uint) *dto.CommentResponseDTO {
	if comment == nil {
		return nil
	}

	username := s.getUsernameForComment(comment.UserID)
	likeCount, isLiked := s.getCommentLikeInfo(comment.ID, userID)

	dto := &dto.CommentResponseDTO{
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
	}

	if comment.Replies != nil && len(comment.Replies) > 0 {
		dto.Replies = s.convertCommentsToDTO(comment.Replies, userID)
	}

	return dto
}

func (s *PostService) getUsernameForComment(userID uint) string {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return "Unknown User"
	}

	if user != nil {
		return user.Username
	}

	return "Deleted User"
}

func (s *PostService) getCommentLikeInfo(commentID, userID uint) (uint, bool) {
	likeCount, err := s.likeRepo.GetLikeCount(model.LikeTypeComment, commentID)
	if err != nil {
		likeCount = 0
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeTypeComment, commentID, userID)
	if err != nil {
		isLiked = false
	}

	return likeCount, isLiked
}

func (s *PostService) LikeEntity(userID uint, input *dto.LikeRequestDTO) error {
	if input == nil {
		return exception.NewBadRequestException("Input cannot be nil", "LIKE_BAD_INPUT", nil)
	}

	// Validate entity exists first
	switch input.EntityType {
	case "challenge":
		challenge, err := s.challengeRepo.GetChallengeByID(input.EntityID, userID)
		if err != nil {
			return err
		}
		if challenge == nil {
			return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.EntityID), "CHALLENGE_NOT_FOUND")
		}
	case "post":
		post, err := s.postRepo.GetPost(input.EntityID)
		if err != nil {
			return err
		}
		if post == nil {
			return exception.NewNotFoundException("Post", fmt.Sprintf("%d", input.EntityID), "POST_NOT_FOUND")
		}
	case "comment":
		comment, err := s.commentRepo.GetComment(input.EntityID)
		if err != nil {
			return err
		}
		if comment == nil {
			return exception.NewNotFoundException("Comment", fmt.Sprintf("%d", input.EntityID), "COMMENT_NOT_FOUND")
		}
	default:
		return exception.NewBadRequestException("Invalid entity type", "INVALID_ENTITY_TYPE", nil)
	}

	// Check if user already liked this entity
	isLiked, err := s.likeRepo.IsUserLiked(model.LikeType(input.EntityType), input.EntityID, userID)
	if err != nil {
		return err
	}
	if isLiked {
		return exception.NewConflictException("Like", "user_id", "USER_ALREADY_LIKED")
	}

	like := &model.Like{
		EntityType: model.LikeType(input.EntityType),
		EntityID:   input.EntityID,
		UserID:     userID,
	}

	return s.likeRepo.CreateLike(like)
}

func (s *PostService) UnlikeEntity(userID uint, input *dto.LikeRequestDTO) error {
	if input == nil {
		return exception.NewBadRequestException("Input cannot be nil", "UNLIKE_BAD_INPUT", nil)
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeType(input.EntityType), input.EntityID, userID)
	if err != nil {
		return err
	}
	if !isLiked {
		return exception.NewBadRequestException("User has not liked this entity", "USER_NOT_LIKED", nil)
	}

	return s.likeRepo.DeleteLike(model.LikeType(input.EntityType), input.EntityID, userID)
}

func (s *PostService) GetPostsByChallenge(challengeID, userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	// Verify challenge exists
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	// Get posts for this challenge
	posts, err := s.postRepo.GetPostsByChallenge(challengeID, offset, limit)
	if err != nil {
		return nil, err
	}

	return s.enrichPostsWithDetails(posts, userID)
}

// Helper method
func (s *PostService) enrichPostsWithDetails(posts []*model.Post, userID uint) ([]*dto.PostResponseDTO, error) {
	postDTOs := make([]*dto.PostResponseDTO, len(posts))

	for i, post := range posts {
		user, err := s.userRepo.GetUserByID(post.UserID)
		if err != nil {
			return nil, err
		}
		username := ""
		if user != nil {
			username = user.Username
		}

		likeCount, _ := s.likeRepo.GetLikeCount(model.LikeTypePost, post.ID)
		commentCount, _ := s.commentRepo.GetCommentCount(model.CommentTypePost, post.ID)
		isLiked, _ := s.likeRepo.IsUserLiked(model.LikeTypePost, post.ID, userID)

		postDTOs[i] = &dto.PostResponseDTO{
			ID:           post.ID,
			UserID:       post.UserID,
			Username:     username,
			Description:  post.Description,
			ChallengeID:  post.ChallengeID,
			Pictures:     post.Pictures,
			LikeCount:    likeCount,
			CommentCount: commentCount,
			IsLiked:      isLiked,
			CreatedAt:    post.CreatedAt,
			UpdatedAt:    post.UpdatedAt,
		}
	}

	return postDTOs, nil
}
