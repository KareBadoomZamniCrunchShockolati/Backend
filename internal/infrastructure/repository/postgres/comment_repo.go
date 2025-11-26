package postgres

import (
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"

	"gorm.io/gorm"
)

type ChallengeCommentRepository struct {
	db *gorm.DB
}

func NewChallengeCommentRepository(db *gorm.DB) *ChallengeCommentRepository {
	return &ChallengeCommentRepository{db: db}
}

func (r *ChallengeCommentRepository) CreateComment(comment *model.ChallengeComment) (*model.ChallengeComment, error) {
	entity := toChallengeCommentEntity(comment)
	result := r.db.Create(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	return toChallengeCommentModelWithID(*entity, comment), nil
}

func (r *ChallengeCommentRepository) GetComment(commentID uint) (*model.ChallengeComment, error) {
	var entity entity.ChallengeCommentEntity
	result := r.db.First(&entity, commentID)
	return handleGetEntity(result.Error, func() *model.ChallengeComment {
		return toChallengeCommentModel(&entity)
	})
}
func (r *ChallengeCommentRepository) GetChallengeComments(challengeID uint, offset, limit int) ([]*model.ChallengeComment, error) {

	var entities []entity.ChallengeCommentEntity
	result := r.db.Where("challenge_id = ?", challengeID).Order("created_at DESC").Offset(int(offset)).Limit(int(limit)).Find(&entities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	comments := make([]*model.ChallengeComment, len(entities))
	for i, e := range entities {
		comments[i] = toChallengeCommentModel(&e)
	}
	return comments, nil
}

func (r *ChallengeCommentRepository) UpdateComment(comment *model.ChallengeComment) (*model.ChallengeComment, error) {
	entity := toChallengeCommentEntity(comment)
	entity.ID = comment.ID

	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	comment.UpdatedAt = entity.UpdatedAt
	return comment, nil
}

func (r *ChallengeCommentRepository) DeleteComment(commentID uint) error {
	result := r.db.Delete(&entity.ChallengeCommentEntity{}, commentID)
	return handleDeleteResult(result.Error, "Comment", commentID)
}

// Helpers
func toChallengeCommentEntity(m *model.ChallengeComment) *entity.ChallengeCommentEntity {
	return &entity.ChallengeCommentEntity{
		ChallengeID: m.ChallengeID,
		UserID:      m.UserID,
		Content:     m.Content,
	}
}

func toChallengeCommentModel(e *entity.ChallengeCommentEntity) *model.ChallengeComment {
	return &model.ChallengeComment{
		ID:          e.ID,
		ChallengeID: e.ChallengeID,
		UserID:      e.UserID,
		Content:     e.Content,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func toChallengeCommentModelWithID(e entity.ChallengeCommentEntity, original *model.ChallengeComment) *model.ChallengeComment {
	original.ID = e.ID
	original.CreatedAt = e.CreatedAt
	original.UpdatedAt = e.UpdatedAt
	return original
}

func handleGetEntity(err error, modelFunc func() *model.ChallengeComment) (*model.ChallengeComment, error) {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(err)
	}
	return modelFunc(), nil
}

func handleDeleteResult(err error, entityType string, entityID uint) error {
	if err != nil {
		return exception.NewRepositoryError(err)
	}
	return nil
}
