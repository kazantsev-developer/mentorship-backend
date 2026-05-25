package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type ProgressHandler struct {
	progressService *services.ProgressService
}

func NewProgressHandler(progressService *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{progressService: progressService}
}

type markViewedRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
}

func (h *ProgressHandler) MarkMaterialViewed(c *gin.Context) {
	studentID := c.GetString("userID")
	var req markViewedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newStatus, err := h.progressService.MarkMaterialViewed(studentID, req.MaterialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": newStatus})
}
