package postgres

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/enum"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/model"
	"challenge-app/internal/infrastructure/repository/postgres/entity"
	"time"

	"gorm.io/gorm"
)

type ChallengeRepository struct {
	db *gorm.DB
}

func NewChallengeRepository(db *gorm.DB) *ChallengeRepository {
	return &ChallengeRepository{db: db}
}

func (r *ChallengeRepository) CreateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, exception.NewRepositoryError(tx.Error)
	}
	chEntity := toChallengeEntity(challenge)
	if err := tx.Create(chEntity).Error; err != nil {
		tx.Rollback()
		return nil, exception.NewRepositoryError(err)
	}
	participant := &entity.ChallengeParticipantEntity{
		ChallengeID: chEntity.ID,
		UserID:      chEntity.CreatorID,
		Status:      1,
	}
	if err := tx.Create(participant).Error; err != nil {
		tx.Rollback()
		return nil, exception.NewRepositoryError(err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	return toChallengeModelWithID(*chEntity, challenge), nil
}

func (r *ChallengeRepository) GetChallengeByID(id uint, userID uint) (*model.ChallengeModel, error) {
	joinedSubQuery := r.db.Select("challenge_id").
		Table("challenge_participants").
		Where("user_id = ? AND status = ?", userID, uint(enum.StatusJoined))

	pendingInviteSubQuery := r.db.Select("challenge_id").
		Table("challenge_invites").
		Where("invitee_id = ? AND status = ?", userID, uint(enum.InvitePending))

	baseQuery := r.db.Table("challenges c").Where(
		r.db.Where("c.id = ? AND c.visibility = ? AND c.is_stopped = false", id, "public").
			Or(r.db.Where("c.id = ? AND c.visibility = ? AND c.is_stopped = false", id, "private")).
			Or(r.db.Where("c.id = ? AND c.visibility = ? AND (c.creator_id = ? OR c.id IN (?) OR c.id IN (?)) AND c.is_stopped = false",
				id, "invite", userID, joinedSubQuery, pendingInviteSubQuery)),
	)

	var entity entity.ChallengeEntity
	if err := baseQuery.First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, exception.NewRepositoryError(err)
	}
	return toChallengeModel(&entity), nil
}

func (r *ChallengeRepository) GetAllChallenges(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	baseQuery := r.db.Table("challenges c").Where("c.is_stopped = false")
	joinedSubQuery := r.db.Select("challenge_id").Table("challenge_participants").Where("user_id = ?", userID)
	query := baseQuery.Where(
		r.db.Where("c.visibility = ?", string(enum.VisibilityPublic)).
			Or(r.db.Where("c.creator_id = ?", userID)).
			Or(r.db.Where("c.id IN (?) AND c.visibility = ?", joinedSubQuery, string(enum.VisibilityInvite))),
	)
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) UpdateChallenge(challenge *model.ChallengeModel) (*model.ChallengeModel, error) {
	entity := toChallengeEntity(challenge)
	entity.ID = challenge.ID
	result := r.db.Save(&entity)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}
	challenge.UpdatedAt = entity.UpdatedAt
	return challenge, nil
}

func (r *ChallengeRepository) DeleteChallenge(id uint) error {
	result := r.db.Delete(&entity.ChallengeEntity{}, id)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

func (r *ChallengeRepository) StopChallenge(id uint) error {
	result := r.db.Model(&entity.ChallengeEntity{}).Where("id = ?", id).Update("is_stopped", true)
	if result.Error != nil {
		return exception.NewRepositoryError(result.Error)
	}
	return nil
}

// List methods
func (r *ChallengeRepository) ListPublicChallenges(offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.db.Table("challenges c").
		Where("c.visibility = ? AND c.is_stopped = false", string(enum.VisibilityPublic)).
		Order("c.created_at DESC")
	return r.executeQuery(query, 0, offset, limit)
}

func (r *ChallengeRepository) getBaseDiscoverableQuery(userID uint) *gorm.DB {
	db := r.db.Table("challenges c").Where("c.is_stopped = false")

	if userID == 0 {
		return db.Where("c.visibility IN ?", []string{"public", "private"})
	}

	joinedSubQuery := r.db.Select("challenge_id").
		Table("challenge_participants").
		Where("user_id = ? AND status = ?", userID, uint(enum.StatusJoined))
	return db.Where(
		r.db.Where("c.visibility IN ?", []string{"public", "private"}).
			Or(r.db.Where("c.visibility = ? AND c.creator_id = ?", "invite", userID)).
			Or(r.db.Where("c.visibility = ? AND c.id IN (?)", "invite", joinedSubQuery)),
	)
}

func (r *ChallengeRepository) ListDiscoverableChallenges(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	return r.executeQuery(r.getBaseDiscoverableQuery(userID).Order("c.created_at DESC"), userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByCreatorID(creatorID, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).Where("c.creator_id = ?", creatorID)
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByCreatorUsername(username string, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).
		Joins("JOIN users u ON u.id = c.creator_id").
		Where("u.username ILIKE ?", "%"+username+"%")
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByCategoryID(categoryID, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).Where("c.category_id = ?", categoryID)
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByCategoryName(name string, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).
		Joins("JOIN challenge_categories cat ON cat.id = c.category_id").
		Where("cat.name ILIKE ?", "%"+name+"%")
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByParticipant(userID uint, offset, limit int) ([]*model.ChallengeModel, error) {
	var chEntities []entity.ChallengeEntity
	result := r.db.Joins("JOIN challenge_participants cp ON cp.challenge_id = challenges.id").
		Where("cp.user_id = ? AND cp.status IN ?", userID, []uint{uint(enum.StatusJoined)}).
		Offset(offset).Limit(limit).Find(&chEntities)
	if result.Error != nil {
		return nil, exception.NewRepositoryError(result.Error)
	}

	challenges := make([]*model.ChallengeModel, len(chEntities))
	for i, e := range chEntities {
		challenges[i] = toChallengeModel(&e)
	}
	return challenges, nil
}

func (r *ChallengeRepository) ListChallengesByParticipantCount(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).Order("c.participant_count DESC")
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesByLikeCount(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).Order("c.like_count DESC")
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListChallengesStartingSoon(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.getBaseDiscoverableQuery(userID).
		Where("c.start_time > NOW()").
		Order("c.start_time ASC")
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) ListTopCreatorsChallenge(offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	query := r.db.Table("challenges c").
		Select("c.*").
		Where("c.visibility = ? AND c.is_stopped = false", "public").
		Order("c.participant_count DESC")
	return r.executeQuery(query, 0, offset, limit)
}

func (r *ChallengeRepository) ListChallengesJoinedByUser(userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	if userID == 0 {
		return []*dto.ChallengePreviewDTO{}, nil
	}
	query := r.db.Table("challenges c").
		Joins("JOIN challenge_participants cp ON cp.challenge_id = c.id").
		Where("cp.user_id = ? AND cp.status = ? AND c.is_stopped = false", userID, uint(enum.StatusJoined))
	return r.executeQuery(query, userID, offset, limit)
}

func (r *ChallengeRepository) SearchChallenges(query string, visibility []enum.ChallengeVisibility, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	visStrings := make([]string, len(visibility))
	for i, v := range visibility {
		visStrings[i] = string(v)
	}

	baseQuery := r.getBaseDiscoverableQuery(userID)
	if len(visStrings) > 0 {
		baseQuery = baseQuery.Where("c.visibility IN ?", visStrings)
	}

	searchQuery := baseQuery.Where("c.title ILIKE ? OR c.description ILIKE ?", "%"+query+"%", "%"+query+"%")
	return r.executeQuery(searchQuery, userID, offset, limit)
}

func (r *ChallengeRepository) executeQuery(query *gorm.DB, userID uint, offset, limit int) ([]*dto.ChallengePreviewDTO, error) {
	var results []struct {
		ID              uint
		Title           string
		Description     string
		Rule            string
		CategoryID      uint
		CreatorID       uint
		Visibility      string
		ImageURL        string
		MaxParticipants uint
		StartTime       time.Time
		EndTime         *time.Time
		Timezone        string
		CreatedAt       time.Time
	}

	err := query.
		Select("c.id, c.title, c.description, c.rule, c.category_id, c.creator_id, c.visibility, c.image_url, c.max_participants, c.start_time, c.end_time, c.timezone, c.created_at").
		Offset(offset).Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}

	if len(results) == 0 {
		return []*dto.ChallengePreviewDTO{}, nil
	}

	var challengeIDs, categoryIDs, creatorIDs []uint
	for _, res := range results {
		challengeIDs = append(challengeIDs, res.ID)
		categoryIDs = append(categoryIDs, res.CategoryID)
		creatorIDs = append(creatorIDs, res.CreatorID)
	}

	categoryMap := r.batchGetCategoryNames(categoryIDs)
	creatorMap := r.batchGetUsernames(creatorIDs)
	likeCounts := r.batchGetLikeCounts(challengeIDs)
	commentCounts := r.batchGetCommentCounts(challengeIDs)
	participantCounts := r.batchGetParticipantCounts(challengeIDs)
	participationMap := r.batchGetUserParticipation(challengeIDs, userID)

	dtos := make([]*dto.ChallengePreviewDTO, len(results))
	for i, res := range results {
		dtos[i] = &dto.ChallengePreviewDTO{
			ID:                  res.ID,
			Title:               res.Title,
			Description:         res.Description,
			Rule:                res.Rule,
			CategoryName:        categoryMap[res.CategoryID],
			CreatorUsername:     creatorMap[res.CreatorID],
			CreatorID:           res.CreatorID,
			Visibility:          enum.ChallengeVisibility(res.Visibility),
			ImageURL:            res.ImageURL,
			MaxParticipants:     res.MaxParticipants,
			CurrentParticipants: int(participantCounts[res.ID]),
			LikeCount:           likeCounts[res.ID],
			CommentCount:        commentCounts[res.ID],
			StartTime:           res.StartTime,
			EndTime:             res.EndTime,
			Timezone:            res.Timezone,
			CreatedAt:           res.CreatedAt,
			IsUserParticipating: participationMap[res.ID],
		}
	}
	return dtos, nil
}

func (r *ChallengeRepository) batchGetCategoryNames(ids []uint) map[uint]string {
	names := make(map[uint]string)
	if len(ids) == 0 {
		return names
	}
	var categories []struct {
		ID   uint
		Name string
	}
	r.db.Table("challenge_categories").Where("id IN ?", ids).Find(&categories)
	for _, c := range categories {
		names[c.ID] = c.Name
	}
	return names
}

func (r *ChallengeRepository) batchGetUsernames(ids []uint) map[uint]string {
	names := make(map[uint]string)
	if len(ids) == 0 {
		return names
	}
	var users []struct {
		ID       uint
		Username string
	}
	r.db.Table("users").Where("id IN ?", ids).Find(&users)
	for _, u := range users {
		names[u.ID] = u.Username
	}
	return names
}

func (r *ChallengeRepository) batchGetLikeCounts(challengeIDs []uint) map[uint]uint {
	counts := make(map[uint]uint)
	if len(challengeIDs) == 0 {
		return counts
	}
	var results []struct {
		ChallengeID uint
		Count       int64
	}
	r.db.Table("likes").
		Select("challenge_id, COUNT(*) as count").
		Where("challenge_id IN ?", challengeIDs).
		Group("challenge_id").
		Find(&results)
	for _, r := range results {
		counts[r.ChallengeID] = uint(r.Count)
	}
	for _, id := range challengeIDs {
		if _, exists := counts[id]; !exists {
			counts[id] = 0
		}
	}
	return counts
}

func (r *ChallengeRepository) batchGetCommentCounts(challengeIDs []uint) map[uint]uint {
	counts := make(map[uint]uint)
	if len(challengeIDs) == 0 {
		return counts
	}
	var results []struct {
		ChallengeID uint
		Count       int64
	}
	r.db.Table("comments").
		Select("challenge_id, COUNT(*) as count").
		Where("challenge_id IN ?", challengeIDs).
		Group("challenge_id").
		Find(&results)
	for _, r := range results {
		counts[r.ChallengeID] = uint(r.Count)
	}
	for _, id := range challengeIDs {
		if _, exists := counts[id]; !exists {
			counts[id] = 0
		}
	}
	return counts
}

func (r *ChallengeRepository) batchGetParticipantCounts(challengeIDs []uint) map[uint]uint {
	counts := make(map[uint]uint)
	if len(challengeIDs) == 0 {
		return counts
	}
	var results []struct {
		ChallengeID uint
		Count       int64
	}
	r.db.Table("challenge_participants").
		Select("challenge_id, COUNT(*) as count").
		Where("challenge_id IN ? AND status = ?", challengeIDs, uint(enum.StatusJoined)).
		Group("challenge_id").
		Find(&results)
	for _, r := range results {
		counts[r.ChallengeID] = uint(r.Count)
	}
	for _, id := range challengeIDs {
		if _, exists := counts[id]; !exists {
			counts[id] = 0
		}
	}
	return counts
}

func (r *ChallengeRepository) batchGetUserParticipation(challengeIDs []uint, userID uint) map[uint]bool {
	if userID == 0 || len(challengeIDs) == 0 {
		return make(map[uint]bool)
	}
	participation := make(map[uint]bool)
	var participantIDs []uint
	r.db.Table("challenge_participants").
		Where("user_id = ? AND challenge_id IN ? AND status = ?", userID, challengeIDs, uint(enum.StatusJoined)).
		Pluck("challenge_id", &participantIDs)
	for _, id := range participantIDs {
		participation[id] = true
	}
	return participation
}

func (r *ChallengeRepository) IsChallengeCreator(challengeID, userID uint) (bool, error) {
	var chEntity entity.ChallengeEntity
	result := r.db.Select("creator_id").Where("id = ?", challengeID).First(&chEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, exception.NewRepositoryError(result.Error)
	}
	return chEntity.CreatorID == userID, nil
}

func (r *ChallengeRepository) GetMutualFollowersInChallenge(userID, challengeID uint) ([]*model.UserModel, error) {
	var userEntities []entity.UserEntity

	err := r.db.Table("users u").
		Joins("INNER JOIN challenge_participants cp ON u.id = cp.user_id").
		Joins("INNER JOIN follows f ON u.id = f.following_id").
		Where("cp.challenge_id = ? AND f.follower_id = ? AND cp.status IN ?",
			challengeID,
			userID,
			[]uint{uint(enum.StatusJoined), uint(enum.StatusPending)}).
		Select("u.id, u.username, u.email, u.bio, u.verified").
		Find(&userEntities).Error

	if err != nil {
		return nil, exception.NewRepositoryError(err)
	}
	userModels := make([]*model.UserModel, len(userEntities))
	for i, u := range userEntities {
		userModels[i] = &model.UserModel{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Bio:      u.Bio,
			Verified: u.Verified,
		}
	}
	return userModels, nil
}

// Helpers
func toChallengeEntity(m *model.ChallengeModel) *entity.ChallengeEntity {
	return &entity.ChallengeEntity{
		Title:           m.Title,
		Description:     m.Description,
		CategoryID:      m.CategoryID,
		CreatorID:       m.CreatorID,
		MaxParticipants: m.MaxParticipants,
		Visibility:      m.Visibility,
		Rule:            m.Rule,
		StartTime:       &m.StartTime,
		EndTime:         m.EndTime,
		Timezone:        m.Timezone,
		ImageURL:        m.ImageURL,
		IsStopped:       m.IsStopped,
		CommentsEnabled: m.CommentsEnabled,
	}
}

func toChallengeModel(e *entity.ChallengeEntity) *model.ChallengeModel {
	return &model.ChallengeModel{
		ID:              e.ID,
		Title:           e.Title,
		Description:     e.Description,
		CategoryID:      e.CategoryID,
		CreatorID:       e.CreatorID,
		MaxParticipants: e.MaxParticipants,
		Visibility:      enum.ChallengeVisibility(e.Visibility),
		Rule:            e.Rule,
		EndTime:         e.EndTime,
		StartTime:       *e.StartTime,
		Timezone:        e.Timezone,
		ImageURL:        e.ImageURL,
		IsStopped:       e.IsStopped,
		CommentsEnabled: e.CommentsEnabled,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func toChallengeModelWithID(e entity.ChallengeEntity, original *model.ChallengeModel) *model.ChallengeModel {
	original.ID = e.ID
	original.CreatedAt = e.CreatedAt
	original.UpdatedAt = e.UpdatedAt
	return original
}
