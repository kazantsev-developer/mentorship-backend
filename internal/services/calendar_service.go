package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type CalendarService struct {
	db *gorm.DB
}

func NewCalendarService(db *gorm.DB) *CalendarService {
	return &CalendarService{db: db}
}

func (s *CalendarService) CreateEvent(buddyID, studentID, title, eventType, description string, start, end *time.Time, reminder bool) (*models.CalendarEvent, error) {
	event := models.CalendarEvent{
		ID:              uuid.New().String(),
		Title:           title,
		StudentID:       studentID,
		BuddyID:         buddyID,
		Type:            eventType,
		StartDatetime:   *start,
		EndDatetime:     end,
		Description:     description,
		ReminderEnabled: reminder,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	err := s.db.Create(&event).Error
	return &event, err
}

func (s *CalendarService) GetEventsForStudent(studentID string) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent
	err := s.db.Where("student_id = ?", studentID).Order("start_datetime asc").Find(&events).Error
	return events, err
}

func (s *CalendarService) GetEventsForBuddy(buddyID string) ([]models.CalendarEvent, error) {
	var events []models.CalendarEvent
	err := s.db.Where("buddy_id = ?", buddyID).Order("start_datetime asc").Find(&events).Error
	return events, err
}

func (s *CalendarService) GetUpcomingEvents(userID string, days int) ([]models.CalendarEvent, error) {
	now := time.Now()
	limit := now.AddDate(0, 0, days)
	var events []models.CalendarEvent
	err := s.db.Where("(student_id = ? OR buddy_id = ?) AND start_datetime >= ? AND start_datetime <= ?", userID, userID, now, limit).
		Order("start_datetime asc").Find(&events).Error
	return events, err
}
