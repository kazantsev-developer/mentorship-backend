package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"gorm.io/gorm"
)

type AssignmentHandler struct {
	service      *services.AssignmentService
	userRepo     *repositories.UserRepository
	progressRepo *repositories.ProgressRepository
	blockRepo    *repositories.BlockRepository
	activityRepo *repositories.ActivityRepository
	db           *gorm.DB
}

func NewAssignmentHandler(
	service *services.AssignmentService,
	userRepo *repositories.UserRepository,
	progressRepo *repositories.ProgressRepository,
	blockRepo *repositories.BlockRepository,
	activityRepo *repositories.ActivityRepository,
	db *gorm.DB,
) *AssignmentHandler {
	return &AssignmentHandler{
		service:      service,
		userRepo:     userRepo,
		progressRepo: progressRepo,
		blockRepo:    blockRepo,
		activityRepo: activityRepo,
		db:           db,
	}
}

func (h *AssignmentHandler) AssignBuddy(c *gin.Context) {
	var req struct {
		StudentID string `json:"student_id" binding:"required"`
		BuddyID   string `json:"buddy_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.AssignBuddyToStudent(req.StudentID, req.BuddyID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func (h *AssignmentHandler) MyStudents(c *gin.Context) {
	userID := c.GetString("userID")
	students, err := h.service.GetStudentsForBuddy(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *AssignmentHandler) GetStudentActivity(c *gin.Context) {
	studentID := c.Param("id")
	activities, err := h.activityRepo.GetStudentActivities(studentID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *AssignmentHandler) GetStudent(c *gin.Context) {
	studentID := c.Param("id")
	student, err := h.userRepo.FindByID(studentID)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func (h *AssignmentHandler) GetStudentRoadmap(c *gin.Context) {
	studentID := c.Param("id")
	buddyID := c.GetString("userID")

	var assignment models.StudentBuddyAssignment
	if err := h.db.Where("buddy_id = ? AND student_id = ?", buddyID, studentID).First(&assignment).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	blocks, err := h.blockRepo.GetAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type BlockWithProgress struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Percent     int    `json:"percent"`
	}
	result := make([]BlockWithProgress, 0, len(blocks))
	for _, block := range blocks {
		status, _ := h.progressRepo.GetBlockStatus(studentID, block.ID)
		percent, _ := h.progressRepo.GetBlockProgressPercent(studentID, block.ID)
		result = append(result, BlockWithProgress{
			ID:          block.ID,
			Title:       block.Title,
			Description: block.Description,
			Status:      status,
			Percent:     percent,
		})
	}
	c.JSON(http.StatusOK, result)
}
