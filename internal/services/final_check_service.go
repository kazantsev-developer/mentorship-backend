package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// FinalCheckService manages final check scheduling and completion
type FinalCheckService struct {
	db *gorm.DB
}

// NewFinalCheckService returns a new FinalCheckService instance
func NewFinalCheckService(db *gorm.DB) *FinalCheckService {
	return &FinalCheckService{db: db}
}

// ScheduleFinalCheck creates a new scheduled final check
func (s *FinalCheckService) ScheduleFinalCheck(studentID, buddyID string, checkType models.FinalCheckType, scheduledAt time.Time) (*models.FinalCheck, error) {
	check := models.FinalCheck{
		ID:          uuid.New().String(),
		StudentID:   studentID,
		BuddyID:     buddyID,
		Type:        checkType,
		Status:      "scheduled",
		ScheduledAt: &scheduledAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.db.Create(&check).Error
	return &check, err
}

// CompleteFinalCheck finalizes a final check as passed or failed
func (s *FinalCheckService) CompleteFinalCheck(checkID string, passed bool) error {
	var check models.FinalCheck
	if err := s.db.First(&check, "id = ?", checkID).Error; err != nil {
		return err
	}
	if check.Status == "completed" || check.Status == "failed" {
		return errors.New("already finalized")
	}
	if passed {
		check.Status = "completed"
	} else {
		check.Status = "failed"
	}
	now := time.Now()
	check.CompletedAt = &now
	check.UpdatedAt = now
	return s.db.Save(&check).Error
}

// GetStudentFinalChecks returns all final checks for a student ordered by most recent
func (s *FinalCheckService) GetStudentFinalChecks(studentID string) ([]models.FinalCheck, error) {
	var checks []models.FinalCheck
	err := s.db.Where("student_id = ?", studentID).Order("created_at desc").Find(&checks).Error
	return checks, err
}
