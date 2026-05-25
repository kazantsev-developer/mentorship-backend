package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type FinalCheckService struct {
	db *gorm.DB
}

func NewFinalCheckService(db *gorm.DB) *FinalCheckService {
	return &FinalCheckService{db: db}
}

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

func (s *FinalCheckService) GetStudentFinalChecks(studentID string) ([]models.FinalCheck, error) {
	var checks []models.FinalCheck
	err := s.db.Where("student_id = ?", studentID).Order("created_at desc").Find(&checks).Error
	return checks, err
}
