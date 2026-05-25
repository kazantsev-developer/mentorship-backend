package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type RoadmapService struct {
	blockRepo *repositories.BlockRepository
}

func NewRoadmapService(blockRepo *repositories.BlockRepository) *RoadmapService {
	return &RoadmapService{blockRepo: blockRepo}
}

// GetFullRoadmap возвращает все блоки с материалами для студента
func (s *RoadmapService) GetFullRoadmap() ([]models.Block, error) {
	return s.blockRepo.GetAllActive()
}
