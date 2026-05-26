package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type AchievementHandler struct {
	achievementRepo *repositories.AchievementRepository
	userID          string
}

func NewAchievementHandler(repo *repositories.AchievementRepository) *AchievementHandler {
	return &AchievementHandler{achievementRepo: repo}
}

func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	userID := c.GetString("userID")
	achievements, err := h.achievementRepo.GetActiveAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Для каждого достижения определяем, получено ли оно пользователем
	result := make([]map[string]interface{}, 0)
	for _, ach := range achievements {
		has, _ := h.achievementRepo.HasUserAchievement(userID, ach.ID)
		result = append(result, map[string]interface{}{
			"id":           ach.ID,
			"title":        ach.Title,
			"description":  ach.Description,
			"reward_bonus": ach.RewardBonus,
			"image_url":    ach.ImageURL,
			"unlocked":     has,
		})
	}
	c.JSON(http.StatusOK, result)
}
