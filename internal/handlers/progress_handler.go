package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

// ProgressHandler handles requests related to learning progress
type ProgressHandler struct {
	progressService *services.ProgressService
}

// NewProgressHandler returns a new ProgressHandler instance
func NewProgressHandler(progressService *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

type markViewedRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
}

// MarkMaterialViewed marks a specific material as viewed by the current student
func (h *ProgressHandler) MarkMaterialViewed(c *gin.Context) {
	studentID := c.GetString("userID")
	var req markViewedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	newStatus, err := h.progressService.MarkMaterialViewed(studentID, req.MaterialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark material as viewed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": newStatus})
}
