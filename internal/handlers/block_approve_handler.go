package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type BlockApproveHandler struct {
	progressRepo *repositories.ProgressRepository
}

func NewBlockApproveHandler(progressRepo *repositories.ProgressRepository) *BlockApproveHandler {
	return &BlockApproveHandler{progressRepo: progressRepo}
}

func (h *BlockApproveHandler) ApproveBlock(c *gin.Context) {
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		BlockID   string `json:"block_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	buddyID := c.GetString("userID")
	if err := h.progressRepo.ConfirmBlock(req.StudentID, req.BlockID, buddyID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
