package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type BlockApproveHandler struct {
	progressRepo *repositories.ProgressRepository
	activityRepo *repositories.ActivityRepository
}

func NewBlockApproveHandler(progressRepo *repositories.ProgressRepository, activityRepo *repositories.ActivityRepository) *BlockApproveHandler {
	return &BlockApproveHandler{progressRepo: progressRepo, activityRepo: activityRepo}
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
	metadata := `{"block_id": "` + req.BlockID + `"}`
	_ = h.activityRepo.CreateActivity(req.StudentID, buddyID, "block_approved", "block", req.BlockID, metadata)

	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
