package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type FinalCheckHandler struct {
	service *services.FinalCheckService
}

func NewFinalCheckHandler(service *services.FinalCheckService) *FinalCheckHandler {
	return &FinalCheckHandler{service: service}
}

func (h *FinalCheckHandler) Schedule(c *gin.Context) {
	var req struct {
		StudentID   string    `json:"student_id" binding:"required"`
		Type        string    `json:"type" binding:"required"`
		ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	buddyID := c.GetString("userID")
	checkType := models.FinalCheckType(req.Type)
	if checkType != models.FinalTechnical && checkType != models.FinalRoast {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}
	check, err := h.service.ScheduleFinalCheck(req.StudentID, buddyID, checkType, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, check)
}

func (h *FinalCheckHandler) Complete(c *gin.Context) {
	var req struct {
		CheckID string `json:"check_id" binding:"required"`
		Passed  bool   `json:"passed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CompleteFinalCheck(req.CheckID, req.Passed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

func (h *FinalCheckHandler) GetForStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		studentID = c.GetString("userID")
	}
	checks, err := h.service.GetStudentFinalChecks(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, checks)
}
