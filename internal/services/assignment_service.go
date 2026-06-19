package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"gorm.io/gorm"
)

// AssignmentService manages buddy-student assignments
type AssignmentService struct {
	db       *gorm.DB
	userRepo *repositories.UserRepository
}

// NewAssignmentService returns a new AssignmentService instance
func NewAssignmentService(db *gorm.DB, userRepo *repositories.UserRepository) *AssignmentService {
	return &AssignmentService{db: db, userRepo: userRepo}
}

// AssignBuddyToStudent links a buddy to a student, replacing any existing assignment
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

// GetBuddyForStudent returns the assigned buddy for a student
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

// GetStudentsForBuddy returns all students assigned to a buddy
func (s *AssignmentService) GetStudentsForBuddy(buddyID string) ([]models.User, error) {
	var assignments []models.StudentBuddyAssignment
	err := s.db.Where("buddy_id = ?", buddyID).Find(&assignments).Error
	if err != nil {
		return nil, err
	}
	students := make([]models.User, 0, len(assignments))
	for _, a := range assignments {
		student, err := s.userRepo.FindByID(a.StudentID)
		if err != nil {
			continue
		}
		if student != nil {
			students = append(students, *student)
		}
	}
	return students, nil
}
