package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"gorm.io/gorm"
)

type AdminOneOnOneHandler struct {
	db           *gorm.DB
	bonusService *services.BonusService
}

func NewAdminOneOnOneHandler(db *gorm.DB, bonusService *services.BonusService) *AdminOneOnOneHandler {
	return &AdminOneOnOneHandler{db: db, bonusService: bonusService}
}

func (h *AdminOneOnOneHandler) ListRequests(c *gin.Context) {
	var requests []models.OneOnOneRequest
	h.db.Order("created_at desc").Find(&requests)

	type RequestWithStudent struct {
		ID           string    `json:"id"`
		StudentID    string    `json:"student_id"`
		StudentName  string    `json:"student_name"`
		StudentBonus int       `json:"student_bonus"`
		Status       string    `json:"status"`
		CreatedAt    time.Time `json:"created_at"`
	}
	result := make([]RequestWithStudent, 0)
	for _, req := range requests {
		var user models.User
		h.db.First(&user, "id = ?", req.StudentID)
		balance, _ := h.bonusService.GetBalance(req.StudentID)
		result = append(result, RequestWithStudent{
			ID:           req.ID,
			StudentID:    req.StudentID,
			StudentName:  user.DisplayName,
			StudentBonus: balance,
			Status:       req.Status,
			CreatedAt:    req.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminOneOnOneHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	var req models.OneOnOneRequest
	if err := h.db.First(&req, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	if req.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request already processed"})
		return
	}
	balance, err := h.bonusService.GetBalance(req.StudentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if balance < 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient bonus balance (need 1000)"})
		return
	}
	// Используем прямой вызов репозитория через экспортируемое поле (если не хотим создавать метод в сервисе)
	// В bonus_service.go поле repo не экспортировано, поэтому добавим публичный метод в bonus_service.go
	// Для простоты создадим метод SpendForOneOnOne в bonus_service.go (см. ниже)
	// Но чтобы не усложнять, добавим метод прямо сейчас:
	_, err = h.bonusService.SpendForOneOnOne(req.StudentID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Status = "approved"
	adminID := c.GetString("userID")
	req.ProcessedBy = &adminID
	now := time.Now()
	req.ProcessedAt = &now
	h.db.Save(&req)
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

func (h *AdminOneOnOneHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	var req models.OneOnOneRequest
	if err := h.db.First(&req, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	if req.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request already processed"})
		return
	}
	req.Status = "rejected"
	adminID := c.GetString("userID")
	req.ProcessedBy = &adminID
	now := time.Now()
	req.ProcessedAt = &now
	h.db.Save(&req)
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}

func (h *AdminOneOnOneHandler) Complete(c *gin.Context) {
	id := c.Param("id")
	var req models.OneOnOneRequest
	if err := h.db.First(&req, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}
	if req.Status != "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request not approved"})
		return
	}
	req.Status = "completed"
	adminID := c.GetString("userID")
	req.ProcessedBy = &adminID
	now := time.Now()
	req.ProcessedAt = &now
	h.db.Save(&req)
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}
