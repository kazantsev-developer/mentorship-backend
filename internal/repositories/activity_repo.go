package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) CreateActivity(userID, actorID, activityType, sourceType, sourceID, metadata string) error {
	var actorIDPtr *string
	if actorID != "" {
		actorIDPtr = &actorID
	}
	event := models.ActivityEvent{
		ID:         uuid.New().String(),
		UserID:     userID,
		ActorID:    actorIDPtr,
		Type:       activityType,
		SourceType: sourceType,
		SourceID:   sourceID,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}
	return r.db.Create(&event).Error
}

func (r *ActivityRepository) GetStudentActivities(studentID string, limit int) ([]models.ActivityEvent, error) {
	var activities []models.ActivityEvent
	err := r.db.Where("user_id = ?", studentID).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities).Error
	return activities, err
}
