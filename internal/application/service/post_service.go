package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/storage"
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo       repository.PostRepository
	commentRepo    repository.CommentRepository
	likeRepo       repository.LikeRepository
	userRepo       repository.UserRepository
	challengeRepo  repository.ChallengeRepository
	objectStorage  storage.ObjectStorage
	tempUploadRepo repository.TempUploadRepository
}

func NewPostService(
	postRepo repository.PostRepository,
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	userRepo repository.UserRepository,
	challengeRepo repository.ChallengeRepository,
	objectStorage storage.ObjectStorage,
	tempUploadRepo repository.TempUploadRepository,
) *PostService {
	return &PostService{
		postRepo:       postRepo,
		commentRepo:    commentRepo,
		likeRepo:       likeRepo,
		userRepo:       userRepo,
		challengeRepo:  challengeRepo,
		objectStorage:  objectStorage,
		tempUploadRepo: tempUploadRepo,
	}
}
func (s *PostService) CreatePost(ctx context.Context, userID uint, input *dto.CreatePostDTO) (*model.Post, error) {
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

	// ✅ NEW RULE: post can have up to 10 images, and images come from TempKeys
	if len(input.TempKeys) > 10 {
		return nil, exception.NewBadRequestException("Maximum 10 pictures allowed", "POST_TOO_MANY_PICTURES", nil)
	}

	// ✅ Create post WITHOUT pictures first
	post := &model.Post{
		UserID:      userID,
		Description: input.Description,
		ChallengeID: input.ChallengeID,
		Pictures:    []string{}, // will be filled after commit
	}

	created, err := s.postRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	// ✅ If no images, return immediately
	if len(input.TempKeys) == 0 {
		return created, nil
	}

	urls, err := s.CommitPostImages(ctx, userID, created.ID, input.TempKeys)
	if err != nil {
		// Best-effort cleanup: if commit failed, try to remove any remaining temp keys
		for _, k := range input.TempKeys {
			_ = s.objectStorage.Delete(ctx, k)
			if s.tempUploadRepo != nil {
				_ = s.tempUploadRepo.Untrack(ctx, k)
			}
		}
		return nil, err
	}

	created.Pictures = urls
	updated, err := s.postRepo.UpdatePost(created)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *PostService) GetPost(postID, userID uint) (*dto.PostResponseDTO, error) {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	user, err := s.userRepo.GetUserByID(post.UserID)
	if err != nil {
		return nil, err
	}
	username := ""
	if user != nil {
		username = user.Username
	}

	likeCount, _ := s.likeRepo.GetLikeCount(model.LikeTypePost, postID)
	commentCount, _ := s.commentRepo.GetCommentCount(model.CommentTypePost, postID)
	isLiked, _ := s.likeRepo.IsUserLiked(model.LikeTypePost, postID, userID)

	return &dto.PostResponseDTO{
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
	}, nil
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
		// Create a copy to avoid modifying the original
		commentCopy := *comment
		commentMap[commentCopy.ID] = &commentCopy

		// If no parent or parent is 0, it's a root comment
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

	// Get username with proper error handling
	username := s.getUsernameForComment(comment.UserID)

	// Get like information with error handling
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

	// Recursively convert replies if they exist
	if comment.Replies != nil && len(comment.Replies) > 0 {
		dto.Replies = s.convertCommentsToDTO(comment.Replies, userID)
	}

	return dto
}

// Helper function to get username with proper error handling
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

	// Validate entity exists first (your existing code)
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


func (s *PostService) PresignPostImages(ctx context.Context, userID uint, req dto.PresignPostImagesRequest) (*dto.PresignPostImagesResponse, error) {
	count := req.Count
	if count < 0 {
		return nil, exception.NewBadRequestException("count must be >= 0", "INVALID_COUNT", nil)
	}
	if count > bootstrap.MaxPostImages {
		return nil, exception.NewBadRequestException("maximum 10 images allowed", "POST_TOO_MANY_IMAGES", nil)
	}
	ct := strings.ToLower(strings.TrimSpace(req.ContentType))
	if ct == "" {
		ct = "image/jpeg"
	}
	if ct != "image/jpeg" && ct != "image/png" {
		return nil, exception.NewBadRequestException("unsupported content_type (only image/jpeg,image/png)", "INVALID_CONTENT_TYPE", nil)
	}
	uploads := make([]dto.PresignedUploadDTO, 0, count)
	for i := 0; i < count; i++ {
		ext := ".jpg"
		if ct == "image/png" {
			ext = ".png"
		}
		key := fmt.Sprintf("tmp/posts/%d/%s%s", userID, uuid.NewString(), ext)
		up, err := s.objectStorage.PresignPut(ctx, key, ct, bootstrap.PostImagesTTL)
		if err != nil {
			return nil, exception.NewInternalServerException("failed to presign upload", "PRESIGN_FAILED", err)
		}
		if s.tempUploadRepo != nil {
			_ = s.tempUploadRepo.Track(ctx, key, userID, up.ExpiresAt)
		}
		uploads = append(uploads, dto.PresignedUploadDTO{
			Key:           up.Key,
			UploadURL:     up.UploadURL,
			Headers:       up.Headers,
			TempPublicURL: s.objectStorage.PublicURL(up.Key),
			ExpiresAt:     up.ExpiresAt.UTC().Format(time.RFC3339),
		})
	}
	return &dto.PresignPostImagesResponse{Uploads: uploads}, nil
}

func (s *PostService) CommitPostImages(ctx context.Context, userID uint, postID uint, tempKeys []string) ([]string, error) {
	if len(tempKeys) == 0 {
		return []string{}, nil
	}
	if len(tempKeys) > bootstrap.MaxPostImages {
		return nil, exception.NewBadRequestException("maximum 10 images allowed", "POST_TOO_MANY_IMAGES", nil)
	}
	publicURLs := make([]string, 0, len(tempKeys))
	for _, tmpKey := range tempKeys {
		if tmpKey == "" {
			continue
		}
		prefix := fmt.Sprintf("tmp/posts/%d/", userID)
		if !strings.HasPrefix(tmpKey, prefix) {
			return nil, exception.NewBadRequestException("invalid temp key", "INVALID_TEMP_KEY", nil)
		}
		if s.tempUploadRepo != nil {
			ok, err := s.tempUploadRepo.VerifyOwnership(ctx, tmpKey, userID)
			if err != nil {
				return nil, exception.NewInternalServerException("failed to verify temp key", "TEMP_VERIFY_FAILED", err)
			}
			if !ok {
				return nil, exception.NewBadRequestException("temp key expired or not owned by user", "TEMP_KEY_NOT_OWNED", nil)
			}
		}
		filename := path.Base(tmpKey)
		dstKey := fmt.Sprintf("posts/%d/%s", postID, filename)
		url, err := s.objectStorage.Copy(ctx, tmpKey, dstKey)
		if err != nil {
			return nil, exception.NewInternalServerException("failed to commit image", "COMMIT_FAILED", err)
		}
		_ = s.objectStorage.Delete(ctx, tmpKey) 
		if s.tempUploadRepo != nil {
			_ = s.tempUploadRepo.Untrack(ctx, tmpKey)
		}
		publicURLs = append(publicURLs, url)
	}

	return publicURLs, nil
}
