package postgres

import (
	"gorm.io/gorm"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"challenge-app/internal/domain/exception"
	"fmt"
)

type ChallengeCommentRepository struct {
	db *gorm.DB
}

func NewChallengeCommentRepo(db *gorm.DB) *ChallengeCommentRepository {
	return &ChallengeCommentRepository{db: db}
}

//Helper functions
func toCommentEntity(m *model.ChallengeComment) *entity.ChallengeCommentEntity {
	return &entity.ChallengeCommentEntity{
		ChallengeID: m.ChallengeID,
		UserID:    m.UserID,
		Content:   m.Content,
	}
}

func toCommentModel(e *entity.ChallengeCommentEntity) *model.ChallengeComment {
	return &model.ChallengeComment{
		ID:         e.ID,
		ChallengeID: e.ChallengeID,
		UserID:    e.UserID,
		Content:   e.Content,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func (r *ChallengeCommentRepository) CreateComment(comment *model.ChallengeComment) error {
	commentEntity := toCommentEntity(comment)
	if err := r.db.Create(commentEntity).Error; err != nil {
		return exception.NewRepositoryError("failed to create comment", err)
	}
	return nil
}

func (r *ChallengeCommentRepository) UpdateComment(comment *model.ChallengeComment) error {
	commentEntity := toCommentEntity(comment)
	if err := r.db.Save(commentEntity).Error; err != nil {
		return exception.NewRepositoryError("failed to update comment", err)
	}
	return nil
}

func (r *ChallengeCommentRepository) DeleteComment(commentID uint) error {
	if err := r.db.Delete(&entity.ChallengeCommentEntity{}, commentID).Error; err != nil {
		return exception.NewRepositoryError(fmt.Sprintf("DeleteComment %d", commentID), err)
	}
	return nil
}

func (r *ChallengeCommentRepository) ListComments(challengeID uint, offset, limit int) ([]*model.ChallengeComment, error) {
	var commentEntities []entity.ChallengeCommentEntity
	if err := r.db.Where("challenge_id = ?", challengeID).Offset(offset).Limit(limit).Find(&commentEntities).Error; err != nil {
		return nil, exception.NewRepositoryError(fmt.Sprintf("ListComments %d", challengeID), err)
	}
	var comments []*model.ChallengeComment
	for _, e := range commentEntities {
		comments = append(comments, toCommentModel(&e))
	}
	return comments, nil
}