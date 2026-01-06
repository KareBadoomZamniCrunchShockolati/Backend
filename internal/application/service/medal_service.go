package service

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	"challenge-app/internal/domain/repository"
	"strings"
)

// least challenges that should be completed in a category to earn each medal type
var medalThresholds = map[string]int{
	"rookie": 1,
	"bronze": 10,
	"silver": 25,
	"gold":   50,
}

type MedalService struct {
	medalRepo     repository.MedalRepository
	userMedalRepo repository.UserMedalRepository
	categoryRepo  repository.CategoryRepository
}

func NewMedalService(medalRepo repository.MedalRepository, userMedalRepo repository.UserMedalRepository, categoryRepo repository.CategoryRepository) *MedalService {
	return &MedalService{medalRepo: medalRepo, userMedalRepo: userMedalRepo, categoryRepo: categoryRepo}
}

// called whenever a user completes a challenge in a category.
func (s *MedalService) OnCategoryCompleted(userID, categoryID uint) error {
	count, err := s.userMedalRepo.IncrementUserCategoryCount(userID, categoryID)
	if err != nil {
		return err
	}

	// check thresholds
	for medalType, threshold := range medalThresholds {
		if count == threshold {
			// award medal
			_, err := s.userMedalRepo.AwardMedal(userID, categoryID, medalType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MedalService) GetUserMedals(userID uint) ([]*dto.MedalDTO, error) {
	ums, err := s.userMedalRepo.GetUserMedals(userID)
	if err != nil {
		return nil, err
	}
	res := make([]*dto.MedalDTO, len(ums))
	for i, m := range ums {
		res[i] = &dto.MedalDTO{Type: m.Type, CategoryID: m.CategoryID, AwardedAt: m.AwardedAt}
	}
	return res, nil
}

func (s *MedalService) GetSelectedMedals(userID uint) ([]string, error) {
	return s.userMedalRepo.GetSelectedMedals(userID)
}

func (s *MedalService) SelectMedals(userID uint, medals []string) error {
	// validate that the user actually owns the medals they are trying to select
	userMedals, err := s.userMedalRepo.GetUserMedals(userID)
	if err != nil {
		return err
	}

	// build set of owned medal identifiers in format "categoryName:medalType"
	owned := make(map[string]bool)
	for _, um := range userMedals {
		cat, err := s.categoryRepo.GetCategoryByID(um.CategoryID)
		if err != nil || cat == nil {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(cat.Name)) + ":" + strings.ToLower(strings.TrimSpace(um.Type))
		owned[key] = true
	}

	// Validate requested selections
	for _, m := range medals {
		norm := strings.ToLower(strings.TrimSpace(m))
		// expect format category:medalType
		parts := strings.SplitN(norm, ":", 2)
		if len(parts) != 2 {
			return exception.NewBadRequestException("INVALID_MEDAL_FORMAT", map[string]any{"medal": m})
		}
		if !owned[norm] {
			return exception.NewBadRequestException("MEDAL_NOT_OWNED", map[string]any{"medal": m})
		}
	}

	return s.userMedalRepo.SetSelectedMedals(userID, medals)
}

// DeleteSelectedMedals removes the listed medals from the user's selected list.
func (s *MedalService) DeleteSelectedMedals(userID uint, medals []string) error {
	// normalize input
	for i, m := range medals {
		medals[i] = strings.ToLower(strings.TrimSpace(m))
	}

	selected, err := s.userMedalRepo.GetSelectedMedals(userID)
	if err != nil {
		return err
	}
	// build set of selected
	selectedSet := make(map[string]bool)
	for _, sMed := range selected {
		selectedSet[strings.ToLower(strings.TrimSpace(sMed))] = true
	}

	// validate all requested medals are currently selected
	for _, m := range medals {
		if !selectedSet[m] {
			return exception.NewBadRequestException("MEDAL_NOT_SELECTED", map[string]any{"medal": m})
		}
	}

	// compute remaining
	remaining := make([]string, 0, len(selected))
	removeSet := make(map[string]bool)
	for _, m := range medals {
		removeSet[m] = true
	}
	for _, sMed := range selected {
		n := strings.ToLower(strings.TrimSpace(sMed))
		if !removeSet[n] {
			remaining = append(remaining, sMed)
		}
	}

	return s.userMedalRepo.SetSelectedMedals(userID, remaining)
}
