package services

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"gorm.io/gorm"
)

// RoadmapService provides aggregated roadmap data for students
type RoadmapService struct {
	blockRepo *repositories.BlockRepository
	db        *gorm.DB
}

// NewRoadmapService returns a new RoadmapService instance
func NewRoadmapService(blockRepo *repositories.BlockRepository, db *gorm.DB) *RoadmapService {
	return &RoadmapService{blockRepo: blockRepo, db: db}
}

// GetFullRoadmap builds a complete roadmap with progress for a student
func (s *RoadmapService) GetFullRoadmap(studentID string) ([]map[string]any, error) {
	blocks, err := s.blockRepo.GetAllActive()
	if err != nil {
		return nil, err
	}
	var result []map[string]any
	for _, block := range blocks {
		var blockProgress models.BlockProgress
		s.db.Where("student_id = ? AND block_id = ?", studentID, block.ID).First(&blockProgress)
		percent := 0
		status := "not_started"
		if blockProgress.ID != "" {
			status = string(blockProgress.Status)
			var totalRequired, viewedRequired int64
			s.db.Model(&models.Material{}).Where("block_id = ? AND is_required = ? AND is_active = ?", block.ID, true, true).Count(&totalRequired)
			s.db.Model(&models.MaterialProgress{}).Where("student_id = ? AND material_id IN (SELECT id FROM materials WHERE block_id = ? AND is_required = ?)", studentID, block.ID, true).Count(&viewedRequired)
			if totalRequired > 0 {
				percent = int(viewedRequired * 100 / totalRequired)
				if percent > 100 {
					percent = 100
				}
			}
		}
		materials := make([]map[string]any, 0)
		for _, m := range block.Materials {
			viewed := false
			var mp models.MaterialProgress
			s.db.Where("student_id = ? AND material_id = ?", studentID, m.ID).First(&mp)
			if mp.ID != "" {
				viewed = true
			}
			materials = append(materials, map[string]any{
				"id":          m.ID,
				"title":       m.Title,
				"description": m.Description,
				"type":        m.Type,
				"url":         m.URL,
				"is_required": m.IsRequired,
				"viewed":      viewed,
			})
		}
		result = append(result, map[string]any{
			"id":          block.ID,
			"title":       block.Title,
			"description": block.Description,
			"progress": map[string]any{
				"status":  status,
				"percent": percent,
			},
			"materials": materials,
		})
	}
	return result, nil
}
