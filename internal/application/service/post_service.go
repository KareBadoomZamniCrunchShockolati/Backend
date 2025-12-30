package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/bootstrap"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/domain/repository"
	"challenge-app/internal/infrastructure/storage"

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
		return nil, exception.NewBadRequestException("POST_CREATE_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if user == nil {
		return nil, exception.NewNotFoundException("User", fmt.Sprintf("%d", userID), "USER_NOT_FOUND")
	}

	if input.ChallengeID != nil {
		challenge, err := s.challengeRepo.GetChallengeByID(*input.ChallengeID, userID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		if challenge == nil {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", *input.ChallengeID), "CHALLENGE_NOT_FOUND")
		}
	}

	tempKeys := make([]string, 0, len(input.Images))
	for _, img := range input.Images {
		if img.TempKey != "" {
			tempKeys = append(tempKeys, img.TempKey)
		}
	}
	if len(tempKeys) > bootstrap.MaxPostImages {
		return nil, exception.NewBadRequestException("POST_TOO_MANY_PICTURES", map[string]any{
			"max":   bootstrap.MaxPostImages,
			"count": len(tempKeys),
		})
	}

	post := &model.Post{
		UserID:      userID,
		Description: input.Description,
		ChallengeID: input.ChallengeID,
		Pictures:    []string{},
	}

	created, err := s.postRepo.CreatePost(post)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	if len(tempKeys) == 0 {
		return created, nil
	}

	urls, err := s.CommitPostImages(ctx, userID, created.ID, tempKeys)
	if err != nil {
		for _, k := range tempKeys {
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
		return nil, exception.NewRepositoryError(err)
	}

	return updated, nil
}

func (s *PostService) GetPost(postID, userID uint) (*dto.PostResponseDTO, error) {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if post == nil {
		return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	u, err := s.userRepo.GetUserByID(post.UserID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	username := ""
	if u != nil {
		username = u.Username
	}

	likeCount, err := s.likeRepo.GetLikeCount(model.LikeTypePost, post.ID)
	if err != nil {
		likeCount = 0
	}
	commentCount, err := s.commentRepo.GetCommentCount(model.CommentTypePost, post.ID)
	if err != nil {
		commentCount = 0
	}
	isLiked, err := s.likeRepo.IsUserLiked(model.LikeTypePost, post.ID, userID)
	if err != nil {
		isLiked = false
	}

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
		return nil, exception.NewRepositoryError(err)
	}
	if post == nil {
		return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	if post.UserID != userID {
		return nil, exception.NewForbiddenException("POST_UPDATE_FORBIDDEN")
	}

	if input.Description != nil {
		post.Description = *input.Description
	}
	updated, err := s.postRepo.UpdatePost(post)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return updated, nil
}

func (s *PostService) DeletePost(postID, userID uint) error {
	post, err := s.postRepo.GetPost(postID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if post == nil {
		return exception.NewNotFoundException("Post", fmt.Sprintf("%d", postID), "POST_NOT_FOUND")
	}

	if post.UserID != userID {
		return exception.NewForbiddenException("POST_DELETE_FORBIDDEN")
	}

	err = s.postRepo.DeletePost(postID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *PostService) GetUserPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	posts, err := s.postRepo.GetPostsByUser(userID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return s.enrichPostsWithDetails(posts, userID)
}

func (s *PostService) GetFeedPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	posts, err := s.postRepo.GetFeedPosts(userID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return s.enrichPostsWithDetails(posts, userID)
}

func (s *PostService) AddComment(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error) {
	if input == nil {
		return nil, exception.NewBadRequestException("COMMENT_ADD_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}

	switch input.EntityType {
	case "challenge":
		ch, err := s.challengeRepo.GetChallengeByID(input.EntityID, userID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		if ch == nil {
			return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.EntityID), "CHALLENGE_NOT_FOUND")
		}
	case "post":
		p, err := s.postRepo.GetPost(input.EntityID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
		}
		if p == nil {
			return nil, exception.NewNotFoundException("Post", fmt.Sprintf("%d", input.EntityID), "POST_NOT_FOUND")
		}
	default:
		return nil, exception.NewBadRequestException("INVALID_ENTITY_TYPE", map[string]any{
			"expected": []string{"challenge", "post"},
			"got":      input.EntityType,
		})
	}

	comment := &model.Comment{
		EntityType: model.CommentType(input.EntityType),
		EntityID:   input.EntityID,
		UserID:     userID,
		Content:    input.Content,
		ParentID:   input.ParentID,
	}

	created, err := s.commentRepo.CreateComment(comment)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return created, nil
}

func (s *PostService) GetComments(entityType string, entityID, userID uint, offset, limit int) ([]*dto.CommentResponseDTO, error) {
	comments, err := s.commentRepo.GetComments(model.CommentType(entityType), entityID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	nested := s.buildCommentTree(comments)
	return s.convertCommentsToDTO(nested, userID), nil
}

func (s *PostService) buildCommentTree(comments []*model.Comment) []*model.Comment {
	commentMap := make(map[uint]*model.Comment)
	var rootComments []*model.Comment

	for _, comment := range comments {
		commentCopy := *comment
		commentMap[commentCopy.ID] = &commentCopy

		if commentCopy.ParentID == nil || *commentCopy.ParentID == 0 {
			rootComments = append(rootComments, &commentCopy)
		}
	}

	for _, comment := range comments {
		if comment.ParentID != nil && *comment.ParentID != 0 {
			if parent, exists := commentMap[*comment.ParentID]; exists {
				parent.AddReply(commentMap[comment.ID])
			}
		}
	}

	return rootComments
}

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

func (s *PostService) convertCommentToDTO(comment *model.Comment, userID uint) *dto.CommentResponseDTO {
	if comment == nil {
		return nil
	}

	username := s.getUsernameForComment(comment.UserID)
	likeCount, isLiked := s.getCommentLikeInfo(comment.ID, userID)

	out := &dto.CommentResponseDTO{
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
		out.Replies = s.convertCommentsToDTO(comment.Replies, userID)
	}

	return out
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
		return exception.NewBadRequestException("LIKE_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}

	switch input.EntityType {
	case "challenge":
		ch, err := s.challengeRepo.GetChallengeByID(input.EntityID, userID)
		if err != nil {
			return exception.NewRepositoryError(err)
		}
		if ch == nil {
			return exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", input.EntityID), "CHALLENGE_NOT_FOUND")
		}
	case "post":
		p, err := s.postRepo.GetPost(input.EntityID)
		if err != nil {
			return exception.NewRepositoryError(err)
		}
		if p == nil {
			return exception.NewNotFoundException("Post", fmt.Sprintf("%d", input.EntityID), "POST_NOT_FOUND")
		}
	case "comment":
		cm, err := s.commentRepo.GetComment(input.EntityID)
		if err != nil {
			return exception.NewRepositoryError(err)
		}
		if cm == nil {
			return exception.NewNotFoundException("Comment", fmt.Sprintf("%d", input.EntityID), "COMMENT_NOT_FOUND")
		}
	default:
		return exception.NewBadRequestException("INVALID_ENTITY_TYPE", map[string]any{
			"expected": []string{"challenge", "post", "comment"},
			"got":      input.EntityType,
		})
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeType(input.EntityType), input.EntityID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if isLiked {
		return exception.NewConflictException("USER_ALREADY_LIKED", "Like", "user_id", fmt.Sprintf("%d", userID))
	}

	like := &model.Like{
		EntityType: model.LikeType(input.EntityType),
		EntityID:   input.EntityID,
		UserID:     userID,
	}
	err = s.likeRepo.CreateLike(like)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *PostService) UnlikeEntity(userID uint, input *dto.LikeRequestDTO) error {
	if input == nil {
		return exception.NewBadRequestException("UNLIKE_BAD_INPUT", map[string]any{
			"reason": "input_is_nil",
		})
	}

	isLiked, err := s.likeRepo.IsUserLiked(model.LikeType(input.EntityType), input.EntityID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	if !isLiked {
		return exception.NewBadRequestException("USER_NOT_LIKED", map[string]any{
			"entity_type": input.EntityType,
			"entity_id":   input.EntityID,
			"user_id":     userID,
		})
	}

	err = s.likeRepo.DeleteLike(model.LikeType(input.EntityType), input.EntityID, userID)
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}

func (s *PostService) GetPostsByChallenge(challengeID, userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) {
	challenge, err := s.challengeRepo.GetChallengeByID(challengeID, userID)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	if challenge == nil {
		return nil, exception.NewNotFoundException("Challenge", fmt.Sprintf("%d", challengeID), "CHALLENGE_NOT_FOUND")
	}

	posts, err := s.postRepo.GetPostsByChallenge(challengeID, offset, limit)
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	return s.enrichPostsWithDetails(posts, userID)
}

func (s *PostService) enrichPostsWithDetails(posts []*model.Post, userID uint) ([]*dto.PostResponseDTO, error) {
	postDTOs := make([]*dto.PostResponseDTO, len(posts))

	for i, post := range posts {
		user, err := s.userRepo.GetUserByID(post.UserID)
		if err != nil {
			return nil, exception.NewRepositoryError(err)
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
		return nil, exception.NewBadRequestException("INVALID_COUNT", map[string]any{
			"min":   0,
			"value": count,
		})
	}
	if count > bootstrap.MaxPostImages {
		return nil, exception.NewBadRequestException("POST_TOO_MANY_IMAGES", map[string]any{
			"max":   bootstrap.MaxPostImages,
			"value": count,
		})
	}
	ct := strings.ToLower(strings.TrimSpace(req.ContentType))
	if ct == "" {
		ct = "image/jpeg"
	}
	if ct != "image/jpeg" && ct != "image/png" {
		return nil, exception.NewBadRequestException("INVALID_CONTENT_TYPE", map[string]any{
			"allowed": []string{"image/jpeg", "image/png"},
			"got":     ct,
		})
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
			return nil, exception.NewInternalServerException("PRESIGN_FAILED", map[string]any{
				"reason": "presign_put_failed",
				"key":    key,
			}, err)
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
		return nil, exception.NewBadRequestException("POST_TOO_MANY_IMAGES", map[string]any{
			"max":   bootstrap.MaxPostImages,
			"value": len(tempKeys),
		})
	}
	publicURLs := make([]string, 0, len(tempKeys))
	for _, tmpKey := range tempKeys {
		if tmpKey == "" {
			continue
		}
		prefix := fmt.Sprintf("tmp/posts/%d/", userID)
		if !strings.HasPrefix(tmpKey, prefix) {
			return nil, exception.NewBadRequestException("INVALID_TEMP_KEY", map[string]any{
				"reason": "invalid_prefix",
				"key":    tmpKey,
				"prefix": prefix,
			})
		}
		if s.tempUploadRepo != nil {
			ok, err := s.tempUploadRepo.VerifyOwnership(ctx, tmpKey, userID)
			if err != nil {
				return nil, exception.NewInternalServerException("TEMP_VERIFY_FAILED", map[string]any{
					"key": tmpKey,
				}, err)
			}
			if !ok {
				return nil, exception.NewBadRequestException("TEMP_KEY_NOT_OWNED", map[string]any{
					"key": tmpKey,
				})
			}
		}
		filename := path.Base(tmpKey)
		dstKey := fmt.Sprintf("posts/%d/%s", postID, filename)
		url, err := s.objectStorage.Copy(ctx, tmpKey, dstKey)
		if err != nil {
			return nil, exception.NewInternalServerException("COMMIT_FAILED", map[string]any{
				"src": tmpKey,
				"dst": dstKey,
			}, err)
		}

		_ = s.objectStorage.Delete(ctx, tmpKey)
		if s.tempUploadRepo != nil {
			_ = s.tempUploadRepo.Untrack(ctx, tmpKey)
		}
		publicURLs = append(publicURLs, url)
	}

	return publicURLs, nil
}

func (s *PostService) UploadPostImages(ctx context.Context, userID uint, files []*multipart.FileHeader) (*dto.UploadPostImagesResponse, error) {
	if len(files) == 0 {
		return &dto.UploadPostImagesResponse{Uploads: []dto.PresignedUploadDTO{}}, nil
	}
	if len(files) > bootstrap.MaxPostImages {
		return nil, exception.NewBadRequestException("POST_TOO_MANY_IMAGES", map[string]any{
			"max":   bootstrap.MaxPostImages,
			"value": len(files),
		})
	}

	out := make([]dto.PresignedUploadDTO, 0, len(files))

	for _, fh := range files {
		if fh == nil {
			continue
		}
		ct := strings.ToLower(strings.TrimSpace(fh.Header.Get("Content-Type")))
		ext := strings.ToLower(path.Ext(fh.Filename))

		if ct == "" {
			if ext == ".png" {
				ct = "image/png"
			} else {
				ct = "image/jpeg"
			}
		}

		if ct != "image/jpeg" && ct != "image/png" {
			return nil, exception.NewBadRequestException("INVALID_CONTENT_TYPE", map[string]any{
				"allowed": []string{"image/jpeg", "image/png"},
				"got":     ct,
			})
		}

		file, err := fh.Open()
		if err != nil {
			return nil, exception.NewBadRequestException("INVALID_REQUEST_BODY", map[string]any{
				"reason": "file_open_failed",
			})
		}

		func() {
			defer file.Close()
		}()

		uploadExt := ".jpg"
		if ct == "image/png" {
			uploadExt = ".png"
		}

		key := fmt.Sprintf("tmp/posts/%d/%s%s", userID, uuid.NewString(), uploadExt)

		publicURL, err := s.objectStorage.Upload(ctx, key, ct, file)
		if err != nil {
			return nil, exception.NewInternalServerException("S3_UPLOAD_FAILED", map[string]any{
				"reason": "upload_failed",
				"key":    key,
			}, err)
		}

		expiresAt := time.Now().Add(bootstrap.PostImagesTTL).UTC()
		if s.tempUploadRepo != nil {
			_ = s.tempUploadRepo.Track(ctx, key, userID, expiresAt)
		}

		out = append(out, dto.PresignedUploadDTO{
			Key:           key,
			TempPublicURL: publicURL,
			ExpiresAt:     expiresAt.Format(time.RFC3339),
		})
	}

	return &dto.UploadPostImagesResponse{Uploads: out}, nil
}
