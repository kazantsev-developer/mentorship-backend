package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type ProgressService struct {
	progressRepo       *repositories.ProgressRepository
	materialRepo       *repositories.MaterialRepository
	achievementService *AchievementService
}

func NewProgressService(progressRepo *repositories.ProgressRepository, materialRepo *repositories.MaterialRepository, achievementService *AchievementService) *ProgressService {
	return &ProgressService{
		progressRepo:       progressRepo,
		materialRepo:       materialRepo,
		achievementService: achievementService,
	}
}

func (s *ProgressService) MarkMaterialViewed(studentID, materialID string) (string, error) {
	// 1. Отметить материал
	if err := s.progressRepo.MarkMaterialViewed(studentID, materialID); err != nil {
		return "", err
	}
	// 2. Получить материал
	material, err := s.materialRepo.GetByID(materialID)
	if err != nil {
		return "", err
	}
	blockID := material.BlockID

	// 3. Получить все материалы блока
	materials, err := s.materialRepo.GetMaterialsByBlockID(blockID)
	if err != nil {
		return "", err
	}
	// 4. Получить все просмотренные материалы студентом
	viewedIDs, err := s.progressRepo.GetViewedMaterialIDs(studentID)
	if err != nil {
		return "", err
	}
	viewedMap := make(map[string]bool)
	for _, id := range viewedIDs {
		viewedMap[id] = true
	}
	// 5. Подсчитать обязательные материалы и просмотренные
	var totalRequired, viewedRequired int
	for _, m := range materials {
		if m.IsRequired {
			totalRequired++
			if viewedMap[m.ID] {
				viewedRequired++
			}
		}
	}
	// 6. Определить новый статус блока
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
	// 7. Обновить BlockProgress
	bp := &models.BlockProgress{
		StudentID: studentID,
		BlockID:   blockID,
		Status:    newStatus,
	}
	if err := s.progressRepo.CreateOrUpdateBlockProgress(bp); err != nil {
		return "", err
	}
	// 8. Выдать достижения за количество просмотренных материалов
	totalViewedAll := len(viewedIDs)
	if totalViewedAll > 0 {
		s.achievementService.CheckAndGrantByMaterialCount(studentID, totalViewedAll)
	}
	return string(newStatus), nil
}
