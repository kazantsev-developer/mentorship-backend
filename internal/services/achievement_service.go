package services

import (
	"encoding/json"

	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type AchievementService struct {
	repo         *repositories.AchievementRepository
	bonusService *BonusService
}

func NewAchievementService(repo *repositories.AchievementRepository, bonusService *BonusService) *AchievementService {
	return &AchievementService{repo: repo, bonusService: bonusService}
}

func (s *AchievementService) CheckAndGrantByMaterialCount(userID string, count int) error {
	achievements, err := s.repo.GetByConditionType("material_count")
	if err != nil {
		return err
	}
	for _, ach := range achievements {
		var params map[string]int
		if err := json.Unmarshal([]byte(ach.ConditionParams), &params); err != nil {
			continue
		}
		required, ok := params["count"]
		if !ok {
			continue
		}
		if count >= required {
			reward, err := s.repo.GrantAchievement(userID, ach.ID)
			if err != nil {
				continue
			}
			if reward > 0 {
				s.bonusService.EarnAchievementBonus(userID, reward, ach.ID)
			}
		}
	}
	return nil
}

func (s *AchievementService) CheckAndGrantByBlockClosed(userID string, blockNumber int) error {
	achievements, err := s.repo.GetByConditionType("block_closed")
	if err != nil {
		return err
	}
	for _, ach := range achievements {
		var params map[string]int
		if err := json.Unmarshal([]byte(ach.ConditionParams), &params); err != nil {
			continue
		}
		requiredBlock, ok := params["block_number"]
		if !ok {
			continue
		}
		if blockNumber == requiredBlock {
			reward, err := s.repo.GrantAchievement(userID, ach.ID)
			if err != nil {
				continue
			}
			if reward > 0 {
				s.bonusService.EarnAchievementBonus(userID, reward, ach.ID)
			}
		}
	}
	return nil
}
