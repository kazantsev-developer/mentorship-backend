package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"gorm.io/gorm"
)

// AdminOneOnOneHandler handles administrative operations on 1:1 requests
type AdminOneOnOneHandler struct {
	db           *gorm.DB
	bonusService *services.BonusService
}

// NewAdminOneOnOneHandler returns a new AdminOneOnOneHandler instance
func NewAdminOneOnOneHandler(db *gorm.DB, bonusService *services.BonusService) *AdminOneOnOneHandler {
	return &AdminOneOnOneHandler{db: db, bonusService: bonusService}
}

// ListRequests returns all 1:1 requests with student details
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
	result := make([]RequestWithStudent, 0, len(requests))
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

// Approve approves a pending 1:1 request after checking bonus balance
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check balance"})
		return
	}
	if balance < 1000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient bonus balance (need 1000)"})
		return
	}
	_, err = h.bonusService.SpendForOneOnOne(req.StudentID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process payment"})
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

// Reject rejects a pending 1:1 request
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

// Complete marks an approved 1:1 request as completed
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
