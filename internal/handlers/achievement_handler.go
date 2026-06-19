// Package handlers contains HTTP handlers for the mentorship platform
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

// AchievementHandler handles achievement-related HTTP requests
type AchievementHandler struct {
	achievementRepo *repositories.AchievementRepository
}

// NewAchievementHandler returns a new AchievementHandler instance
func NewAchievementHandler(repo *repositories.AchievementRepository) *AchievementHandler {
	return &AchievementHandler{achievementRepo: repo}
}

// GetAchievements retrieves the list of achievements with unlock status for the current user
func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	userID := c.GetString("userID")
	achievements, err := h.achievementRepo.GetActiveAchievements()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch achievements"})
		return
	}

	result := make([]map[string]any, 0, len(achievements))
	for _, ach := range achievements {
		has, _ := h.achievementRepo.HasUserAchievement(userID, ach.ID)
		result = append(result, map[string]any{
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
