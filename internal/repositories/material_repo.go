package repositories

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) GetByID(id string) (*models.Material, error) {
	var material models.Material
	err := r.db.Where("id = ?", id).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *MaterialRepository) GetMaterialsByBlockID(blockID string) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.Where("block_id = ? AND is_active = ?", blockID, true).Find(&materials).Error
	return materials, err
}
