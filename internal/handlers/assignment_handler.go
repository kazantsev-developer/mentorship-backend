package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

// AssignmentHandler manages buddy-student assignments
type AssignmentHandler struct {
	service *services.AssignmentService
}

// NewAssignmentHandler returns a new AssignmentHandler instance
func NewAssignmentHandler(service *services.AssignmentService) *AssignmentHandler {
	return &AssignmentHandler{service: service}
}

// AssignBuddy links a buddy to a student
func (h *AssignmentHandler) AssignBuddy(c *gin.Context) {
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		BuddyID   string `json:"buddy_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.service.AssignBuddyToStudent(req.StudentID, req.BuddyID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to assign buddy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

// MyStudents returns the list of students assigned to the current buddy
func (h *AssignmentHandler) MyStudents(c *gin.Context) {
	userID := c.GetString("userID")
	students, err := h.service.GetStudentsForBuddy(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch students"})
		return
	}
	c.JSON(http.StatusOK, students)
}
