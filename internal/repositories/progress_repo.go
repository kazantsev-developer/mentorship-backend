package repositories

import (
	"time"

	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// ProgressRepository handles database operations for student progress
type ProgressRepository struct {
	db *gorm.DB
}

// NewProgressRepository returns a new ProgressRepository instance
func NewProgressRepository(db *gorm.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// MarkMaterialViewed records a material as viewed by a student if not already recorded
func (r *ProgressRepository) MarkMaterialViewed(studentID, materialID string) error {
	var existing models.MaterialProgress
	err := r.db.Where("student_id = ? AND material_id = ?", studentID, materialID).First(&existing).Error
	if err == nil {
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

// GetViewedMaterialIDs returns IDs of materials viewed by a student
func (r *ProgressRepository) GetViewedMaterialIDs(studentID string) ([]string, error) {
	var ids []string
	err := r.db.Model(&models.MaterialProgress{}).
		Where("student_id = ?", studentID).
		Pluck("material_id", &ids).Error
	return ids, err
}

// GetBlockProgress returns the block progress for a student
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

// CreateOrUpdateBlockProgress creates or updates a block progress record
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

// ConfirmBlock approves a block for a student
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
