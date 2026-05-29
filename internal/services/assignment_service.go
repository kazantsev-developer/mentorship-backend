package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"gorm.io/gorm"
)

type AssignmentService struct {
	db           *gorm.DB
	userRepo     *repositories.UserRepository
	progressRepo *repositories.ProgressRepository
}

func NewAssignmentService(
	db *gorm.DB,
	userRepo *repositories.UserRepository,
	progressRepo *repositories.ProgressRepository,
) *AssignmentService {
	return &AssignmentService{
		db:           db,
		userRepo:     userRepo,
		progressRepo: progressRepo,
	}
}

func (s *AssignmentService) AssignBuddyToStudent(studentID, buddyID string) error {
	student, err := s.userRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}
	buddy, err := s.userRepo.FindByID(buddyID)
	if err != nil {
		return err
	}
	if buddy == nil {
		return errors.New("buddy not found")
	}
	if err := s.db.Where("student_id = ?", studentID).Delete(&models.StudentBuddyAssignment{}).Error; err != nil {
		return err
	}
	assignment := models.StudentBuddyAssignment{
		ID:        uuid.New().String(),
		StudentID: studentID,
		BuddyID:   buddyID,
		CreatedAt: time.Now(),
	}
	return s.db.Create(&assignment).Error
}

func (s *AssignmentService) GetBuddyForStudent(studentID string) (*models.User, error) {
	var assignment models.StudentBuddyAssignment
	err := s.db.Where("student_id = ?", studentID).First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return s.userRepo.FindByID(assignment.BuddyID)
}

func (s *AssignmentService) GetStudentsForBuddy(buddyID string) ([]models.BuddyStudentResponse, error) {
	var assignments []models.StudentBuddyAssignment
	err := s.db.Where("buddy_id = ?", buddyID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	students := make([]models.BuddyStudentResponse, 0, len(assignments))
	for _, a := range assignments {
		student, err := s.userRepo.FindByID(a.StudentID)
		if err != nil || student == nil {
			continue
		}
		currentBlock, err := s.progressRepo.GetStudentCurrentBlock(student.ID)
		if err != nil {
			continue
		}
		currentBlockTitle := "Не назначен"
		progressPercent := 0
		status := "not_started"
		if currentBlock != nil {
			currentBlockTitle = currentBlock.Title
			progressPercent, _ = s.progressRepo.GetBlockProgressPercent(student.ID, currentBlock.ID)
			status, _ = s.progressRepo.GetBlockStatus(student.ID, currentBlock.ID)
		}
		daysInactive, _ := s.progressRepo.GetLastActivityDays(student.ID)

		students = append(students, models.BuddyStudentResponse{
			User:              *student,
			CurrentBlockTitle: currentBlockTitle,
			ProgressPercent:   progressPercent,
			Status:            status,
			DaysInactive:      daysInactive,
		})
	}
	return students, nil
}
