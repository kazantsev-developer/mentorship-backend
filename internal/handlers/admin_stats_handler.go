package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type AdminStatsHandler struct {
	db *gorm.DB
}

func NewAdminStatsHandler(db *gorm.DB) *AdminStatsHandler {
	return &AdminStatsHandler{db: db}
}

func (h *AdminStatsHandler) GetStats(c *gin.Context) {
	var totalUsers int64
	var studentsCount int64
	var buddiesCount int64
	var activeOneOnOne int64
	var totalAchievementsIssued int64

	h.db.Model(&models.User{}).Where("is_deleted = ?", false).Count(&totalUsers)
	h.db.Model(&models.User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").Where("user_roles.role = ? AND users.is_deleted = ?", "student", false).Count(&studentsCount)
	h.db.Model(&models.User{}).Joins("JOIN user_roles ON user_roles.user_id = users.id").Where("user_roles.role = ? AND users.is_deleted = ?", "buddy", false).Count(&buddiesCount)
	h.db.Model(&models.OneOnOneRequest{}).Where("status = ?", "pending").Count(&activeOneOnOne)
	h.db.Model(&models.UserAchievement{}).Count(&totalAchievementsIssued)

	c.JSON(http.StatusOK, gin.H{
		"total_users":               totalUsers,
		"students_count":            studentsCount,
		"buddies_count":             buddiesCount,
		"active_one_on_one":         activeOneOnOne,
		"total_achievements_issued": totalAchievementsIssued,
	})
}
