package serviceinterface

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/model"
)

type PostServicer interface {
	CreatePost(userID uint, input *dto.CreatePostDTO) (*model.Post, error)
	GetPost(postID, userID uint) (*dto.PostResponseDTO, error)
	UpdatePost(postID, userID uint, input *dto.UpdatePostDTO) (*model.Post, error)
	DeletePost(postID, userID uint) error
	GetUserPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error)
	GetFeedPosts(userID uint, offset, limit int) ([]*dto.PostResponseDTO, error)
	GetPostsByChallenge(challengeID, userID uint, offset, limit int) ([]*dto.PostResponseDTO, error) // Add this
	AddComment(userID uint, input *dto.CommentRequestDTO) (*model.Comment, error)
	GetComments(entityType string, entityID, userID uint, offset, limit int) ([]*dto.CommentResponseDTO, error)
	LikeEntity(userID uint, input *dto.LikeRequestDTO) error
	UnlikeEntity(userID uint, input *dto.LikeRequestDTO) error
}
