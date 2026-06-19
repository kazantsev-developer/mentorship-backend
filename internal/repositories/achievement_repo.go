// Package repositories provides data access operations for the mentorship platform
package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

// AchievementRepository handles database operations for achievements
type AchievementRepository struct {
	db *gorm.DB
}

// NewAchievementRepository returns a new AchievementRepository instance
func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// GetActiveAchievements returns all active achievements ordered by sort order
func (r *AchievementRepository) GetActiveAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&achievements).Error
	return achievements, err
}

// GetByConditionType returns active achievements matching a condition type
func (r *AchievementRepository) GetByConditionType(condType string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Where("condition_type = ? AND is_active = ?", condType, true).Find(&achievements).Error
	return achievements, err
}

// HasUserAchievement checks if a user already has a specific achievement
func (r *AchievementRepository) HasUserAchievement(userID, achievementID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAchievement{}).Where("user_id = ? AND achievement_id = ?", userID, achievementID).Count(&count).Error
	return count > 0, err
}

// GrantAchievement awards an achievement to a user and returns the bonus amount
func (r *AchievementRepository) GrantAchievement(userID, achievementID string) (int, error) {
	has, err := r.HasUserAchievement(userID, achievementID)
	if err != nil {
		return 0, err
	}
	if has {
		return 0, nil
	}
	var ach models.Achievement
	if err := r.db.First(&ach, "id = ?", achievementID).Error; err != nil {
		return 0, err
	}
	ua := models.UserAchievement{
		ID:            uuid.New().String(),
		UserID:        userID,
		AchievementID: achievementID,
		ReceivedAt:    time.Now(),
	}
	if err := r.db.Create(&ua).Error; err != nil {
		return 0, err
	}
	return ach.RewardBonus, nil
}
