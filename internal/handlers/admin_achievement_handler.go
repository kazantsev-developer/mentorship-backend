package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type AdminAchievementHandler struct {
	db *gorm.DB
}

func NewAdminAchievementHandler(db *gorm.DB) *AdminAchievementHandler {
	return &AdminAchievementHandler{db: db}
}

func (h *AdminAchievementHandler) ListAchievements(c *gin.Context) {
	var achievements []models.Achievement
	h.db.Order("sort_order").Find(&achievements)
	c.JSON(http.StatusOK, achievements)
}

func (h *AdminAchievementHandler) CreateAchievement(c *gin.Context) {
	var ach models.Achievement
	if err := c.ShouldBindJSON(&ach); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ach.ID = uuid.New().String()
	ach.CreatedAt = time.Now()
	ach.UpdatedAt = time.Now()
	if err := h.db.Create(&ach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ach)
}

func (h *AdminAchievementHandler) UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var ach models.Achievement
	if err := h.db.First(&ach, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "achievement not found"})
		return
	}
	var updates models.Achievement
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ach.Title = updates.Title
	ach.Description = updates.Description
	ach.RewardBonus = updates.RewardBonus
	ach.ImageURL = updates.ImageURL
	ach.ConditionType = updates.ConditionType
	ach.ConditionParams = updates.ConditionParams
	ach.IsActive = updates.IsActive
	ach.SortOrder = updates.SortOrder
	ach.UpdatedAt = time.Now()
	if err := h.db.Save(&ach).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ach)
}

func (h *AdminAchievementHandler) DeleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&models.Achievement{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminAchievementHandler) ToggleAchievementStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&models.Achievement{}).Where("id = ?", id).Update("is_active", req.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *AdminAchievementHandler) GetAchievementUsers(c *gin.Context) {
	id := c.Param("id")
	var userAchievements []struct {
		UserID     string    `json:"user_id"`
		UserName   string    `json:"display_name"`
		ReceivedAt time.Time `json:"received_at"`
	}
	err := h.db.Table("user_achievements").
		Select("user_achievements.user_id, users.display_name, user_achievements.received_at").
		Joins("JOIN users ON users.id = user_achievements.user_id").
		Where("user_achievements.achievement_id = ?", id).
		Scan(&userAchievements).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userAchievements)
}
