package repositories

import (
	"time"

	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type ProgressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) MarkMaterialViewed(studentID, materialID string) error {
	now := time.Now()
	mp := models.MaterialProgress{
		StudentID:  studentID,
		MaterialID: materialID,
		ViewedAt:   &now,
	}
	return r.db.Where("student_id = ? AND material_id = ?", studentID, materialID).
		Assign(models.MaterialProgress{ViewedAt: &now}).
		FirstOrCreate(&mp).Error
}

func (r *ProgressRepository) GetViewedMaterialIDs(studentID string) ([]string, error) {
	var ids []string
	err := r.db.Model(&models.MaterialProgress{}).
		Where("student_id = ? AND viewed_at IS NOT NULL", studentID).
		Pluck("material_id", &ids).Error
	return ids, err
}

func (r *ProgressRepository) CreateOrUpdateBlockProgress(bp *models.BlockProgress) error {
	var existing models.BlockProgress
	err := r.db.Where("student_id = ? AND block_id = ?", bp.StudentID, bp.BlockID).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Создаём новую запись
			bp.CreatedAt = time.Now()
			bp.UpdatedAt = time.Now()
			return r.db.Create(bp).Error
		}
		return err
	}
	existing.Status = bp.Status
	existing.UpdatedAt = time.Now()
	if bp.ApprovedBy != nil {
		existing.ApprovedBy = bp.ApprovedBy
	}
	if bp.ApprovedAt != nil {
		existing.ApprovedAt = bp.ApprovedAt
	}
	return r.db.Save(&existing).Error
}

func (r *ProgressRepository) GetStudentCurrentBlock(studentID string) (*models.Block, error) {
	var blocks []models.Block
	if err := r.db.Where("is_active = ?", true).Order("sort_order").Find(&blocks).Error; err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return nil, nil
	}

	blockStatuses := make(map[string]string)
	var blockProgress []models.BlockProgress
	if err := r.db.Where("student_id = ?", studentID).Find(&blockProgress).Error; err != nil {
		return nil, err
	}
	for _, bp := range blockProgress {
		blockStatuses[bp.BlockID] = string(bp.Status)
	}

	for _, b := range blocks {
		status, ok := blockStatuses[b.ID]
		if !ok || status != "approved" {
			return &b, nil
		}
	}
	return &blocks[len(blocks)-1], nil
}

func (r *ProgressRepository) GetBlockProgressPercent(studentID, blockID string) (int, error) {
	var totalRequired int64
	if err := r.db.Model(&models.Material{}).
		Where("block_id = ? AND is_required = ? AND is_active = ?", blockID, true, true).
		Count(&totalRequired).Error; err != nil {
		return 0, err
	}
	if totalRequired == 0 {
		return 100, nil
	}

	var viewedRequired int64
	if err := r.db.Model(&models.MaterialProgress{}).
		Joins("JOIN materials ON materials.id = material_progresses.material_id").
		Where("material_progresses.student_id = ? AND materials.block_id = ? AND materials.is_required = ? AND material_progresses.viewed_at IS NOT NULL", studentID, blockID, true).
		Count(&viewedRequired).Error; err != nil {
		return 0, err
	}
	return int(viewedRequired * 100 / totalRequired), nil
}

func (r *ProgressRepository) GetBlockStatus(studentID, blockID string) (string, error) {
	var bp models.BlockProgress
	err := r.db.Where("student_id = ? AND block_id = ?", studentID, blockID).First(&bp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "not_started", nil
		}
		return "", err
	}
	return string(bp.Status), nil
}

func (r *ProgressRepository) GetLastActivityDays(studentID string) (int, error) {
	var lastViewedAt *time.Time
	err := r.db.Model(&models.MaterialProgress{}).
		Where("student_id = ? AND viewed_at IS NOT NULL", studentID).
		Order("viewed_at DESC").
		Limit(1).
		Pluck("viewed_at", &lastViewedAt).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	if lastViewedAt == nil {
		var user models.User
		if err := r.db.Select("learning_started_at").Where("id = ?", studentID).First(&user).Error; err != nil {
			return 0, err
		}
		if user.LearningStartedAt == nil {
			return 0, nil
		}
		lastViewedAt = user.LearningStartedAt
	}
	days := int(time.Since(*lastViewedAt).Hours() / 24)
	if days < 0 {
		days = 0
	}
	return days, nil
}

func (r *ProgressRepository) ConfirmBlock(studentID, blockID, buddyID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var bp models.BlockProgress
		err := tx.Where("student_id = ? AND block_id = ?", studentID, blockID).First(&bp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		now := time.Now()
		if err == gorm.ErrRecordNotFound {
			bp = models.BlockProgress{
				StudentID:  studentID,
				BlockID:    blockID,
				Status:     models.BlockStatusApproved,
				ApprovedBy: &buddyID,
				ApprovedAt: &now,
			}
			return tx.Create(&bp).Error
		}
		bp.Status = models.BlockStatusApproved
		bp.ApprovedBy = &buddyID
		bp.ApprovedAt = &now
		return tx.Save(&bp).Error
	})
}
