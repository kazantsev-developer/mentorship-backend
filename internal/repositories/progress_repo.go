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

// MarkMaterialViewed отмечает материал как просмотренный (если ещё не отмечен)
func (r *ProgressRepository) MarkMaterialViewed(studentID, materialID string) error {
	var existing models.MaterialProgress
	err := r.db.Where("student_id = ? AND material_id = ?", studentID, materialID).First(&existing).Error
	if err == nil {
		// уже отмечен
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	mp := models.MaterialProgress{
		StudentID:  studentID,
		MaterialID: materialID,
		ViewedAt:   time.Now(),
	}
	return r.db.Create(&mp).Error
}

// GetViewedMaterialIDs возвращает список ID материалов, просмотренных студентом
func (r *ProgressRepository) GetViewedMaterialIDs(studentID string) ([]string, error) {
	var ids []string
	err := r.db.Model(&models.MaterialProgress{}).
		Where("student_id = ?", studentID).
		Pluck("material_id", &ids).Error
	return ids, err
}

// GetBlockProgress возвращает текущий статус блока для студента
func (r *ProgressRepository) GetBlockProgress(studentID, blockID string) (*models.BlockProgress, error) {
	var bp models.BlockProgress
	err := r.db.Where("student_id = ? AND block_id = ?", studentID, blockID).First(&bp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bp, nil
}

// CreateOrUpdateBlockProgress создаёт или обновляет статус блока
func (r *ProgressRepository) CreateOrUpdateBlockProgress(bp *models.BlockProgress) error {
	var existing models.BlockProgress
	err := r.db.Where("student_id = ? AND block_id = ?", bp.StudentID, bp.BlockID).First(&existing).Error
	if err == nil {
		bp.ID = existing.ID
		bp.CreatedAt = existing.CreatedAt
		return r.db.Save(bp).Error
	}
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(bp).Error
	}
	return err
}

// ConfirmBlock подтверждает блок (buddy/admin)
func (r *ProgressRepository) ConfirmBlock(studentID, blockID, approvedBy string) error {
	bp := models.BlockProgress{
		StudentID:  studentID,
		BlockID:    blockID,
		Status:     models.BlockApproved,
		ApprovedBy: &approvedBy,
		ApprovedAt: timePtr(time.Now()),
	}
	return r.CreateOrUpdateBlockProgress(&bp)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
