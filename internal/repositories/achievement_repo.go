package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) GetActiveAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Where("is_active = ?", true).Order("sort_order asc").Find(&achievements).Error
	return achievements, err
}

func (r *AchievementRepository) GetByConditionType(condType string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := r.db.Where("condition_type = ? AND is_active = ?", condType, true).Find(&achievements).Error
	return achievements, err
}

func (r *AchievementRepository) HasUserAchievement(userID, achievementID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserAchievement{}).Where("user_id = ? AND achievement_id = ?", userID, achievementID).Count(&count).Error
	return count > 0, err
}

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
