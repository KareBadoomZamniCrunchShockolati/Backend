package serviceinterface

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/model"
	"context"
	"mime/multipart"
)

type PostServicer interface {
	CreatePost(ctx context.Context, userID uint, input *dto.CreatePostDTO) (*model.Post, error)
	GetPost(postID, userID uint) (*dto.PostResponseDTO, error)
	UpdatePost(postID, userID uint, input *dto.UpdatePostDTO) (*model.Post, error)
	DeletePost(postID, userID uint) error
	GetUserPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error)
	GetFeedPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error)
	GetPostsByChallenge(challengeID, userID uint, offset, limit int) ([]*dto.PostResponseDTO, error)

	// Polymorphic methods for posts
	AddComment(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error)
	GetComments(entityType string, entityID, userID uint, offset, limit int) ([]*dto.CommentResponseDTO, error)
	LikeEntity(userID uint, input *dto.LikeRequestDTO) error
	UnlikeEntity(userID uint, input *dto.LikeRequestDTO) error
	PresignPostImages(ctx context.Context, userID uint, req dto.PresignPostImagesRequest) (*dto.PresignPostImagesResponse, error)
	CommitPostImages(ctx context.Context, userID uint, postID uint, tempKeys []string) ([]string, error)
	UploadPostImages(ctx context.Context, userID uint, files []*multipart.FileHeader) (*dto.UploadPostImagesResponse, error)
}
