package repositories

import (
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type BlockRepository struct {
	db *gorm.DB
}

func NewBlockRepository(db *gorm.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

// GetAllActive возвращает все активные блоки с материалами, отсортированные по SortOrder
func (r *BlockRepository) GetAllActive() ([]models.Block, error) {
	var blocks []models.Block
	err := r.db.Where("is_active = ? AND deleted_at IS NULL", true).
		Order("sort_order asc").
		Preload("Materials", "is_active = ? AND deleted_at IS NULL", true).
		Find(&blocks).Error
	return blocks, err
}

// GetByID возвращает блок по ID с его материалами
func (r *BlockRepository) GetByID(id string) (*models.Block, error) {
	var block models.Block
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Materials", "is_active = ? AND deleted_at IS NULL", true).
		First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}
