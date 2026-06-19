package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// InterviewService manages interview records for students and buddies
type InterviewService struct {
	db *gorm.DB
}

// NewInterviewService returns a new InterviewService instance
func NewInterviewService(db *gorm.DB) *InterviewService {
	return &InterviewService{db: db}
}

// InterviewInput contains data for creating an interview
type InterviewInput struct {
	URL      string
	Company  string
	Position string
	Grade    string
	Stack    string
	Date     time.Time
}

// CreateReal records a real interview for a student
func (s *InterviewService) CreateReal(studentID string, input InterviewInput) (*models.Interview, error) {
	interview := models.Interview{
		ID:        uuid.New().String(),
		Type:      models.InterviewReal,
		StudentID: studentID,
		URL:       input.URL,
		Company:   input.Company,
		Position:  input.Position,
		Grade:     input.Grade,
		Stack:     input.Stack,
		Date:      input.Date,
		Status:    "submitted",
		Result:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.db.Create(&interview).Error
	return &interview, err
}

// CreateMock schedules a mock interview for a student with a buddy
func (s *InterviewService) CreateMock(buddyID, studentID string, input InterviewInput) (*models.Interview, error) {
	buddyIDPtr := &buddyID
	interview := models.Interview{
		ID:        uuid.New().String(),
		Type:      models.InterviewMock,
		StudentID: studentID,
		BuddyID:   buddyIDPtr,
		URL:       input.URL,
		Company:   input.Company,
		Position:  input.Position,
		Grade:     input.Grade,
		Stack:     input.Stack,
		Date:      input.Date,
		Status:    "scheduled",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.db.Create(&interview).Error
	return &interview, err
}

// CompleteMock finalizes a mock interview with feedback
func (s *InterviewService) CompleteMock(interviewID string, feedback string) error {
	var interview models.Interview
	if err := s.db.First(&interview, "id = ?", interviewID).Error; err != nil {
		return err
	}
	if interview.Type != models.InterviewMock {
		return errors.New("not a mock interview")
	}
	interview.Status = "completed"
	interview.Feedback = feedback
	interview.UpdatedAt = time.Now()
	return s.db.Save(&interview).Error
}

// GetUserInterviews returns interviews for a user as a student or as a buddy
func (s *InterviewService) GetUserInterviews(userID string, isBuddy bool) ([]models.Interview, error) {
	var interviews []models.Interview
	query := s.db.Where("student_id = ?", userID)
	if isBuddy {
		query = s.db.Where("buddy_id = ?", userID)
	}
	err := query.Order("date desc").Find(&interviews).Error
	return interviews, err
}

// GetAllRealInterviews returns all real interviews across all students
func (s *InterviewService) GetAllRealInterviews() ([]models.Interview, error) {
	var interviews []models.Interview
	err := s.db.Where("type = ?", models.InterviewReal).Order("date desc").Find(&interviews).Error
	return interviews, err
}
