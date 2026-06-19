package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

// ProgressService tracks student progress through roadmap blocks and materials
type ProgressService struct {
	progressRepo       *repositories.ProgressRepository
	materialRepo       *repositories.MaterialRepository
	achievementService *AchievementService
}

// NewProgressService returns a new ProgressService instance
func NewProgressService(progressRepo *repositories.ProgressRepository, materialRepo *repositories.MaterialRepository, achievementService *AchievementService) *ProgressService {
	return &ProgressService{
		progressRepo:       progressRepo,
		materialRepo:       materialRepo,
		achievementService: achievementService,
	}
}

// MarkMaterialViewed records a material as viewed and updates block progress
func (s *ProgressService) MarkMaterialViewed(studentID, materialID string) (string, error) {
	if err := s.progressRepo.MarkMaterialViewed(studentID, materialID); err != nil {
		return "", err
	}
	material, err := s.materialRepo.GetByID(materialID)
	if err != nil {
		return "", err
	}
	blockID := material.BlockID

	materials, err := s.materialRepo.GetMaterialsByBlockID(blockID)
	if err != nil {
		return "", err
	}
	viewedIDs, err := s.progressRepo.GetViewedMaterialIDs(studentID)
	if err != nil {
		return "", err
	}
	viewedMap := make(map[string]bool)
	for _, id := range viewedIDs {
		viewedMap[id] = true
	}
	var totalRequired, viewedRequired int
	for _, m := range materials {
		if m.IsRequired {
			totalRequired++
			if viewedMap[m.ID] {
				viewedRequired++
			}
		}
	}
	var newStatus models.BlockStatus
	if totalRequired == 0 {
		newStatus = models.BlockInProgress
	} else if viewedRequired == totalRequired {
		newStatus = models.BlockWaitingBuddyConfirmation
	} else if viewedRequired > 0 {
		newStatus = models.BlockInProgress
	} else {
		newStatus = models.BlockNotStarted
	}
	bp := &models.BlockProgress{
		StudentID: studentID,
		BlockID:   blockID,
		Status:    newStatus,
	}
	if err := s.progressRepo.CreateOrUpdateBlockProgress(bp); err != nil {
		return "", err
	}
	totalViewedAll := len(viewedIDs)
	if totalViewedAll > 0 {
		s.achievementService.CheckAndGrantByMaterialCount(studentID, totalViewedAll)
	}
	return string(newStatus), nil
}
