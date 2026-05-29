package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type InterviewHandler struct {
	service      *services.InterviewService
	activityRepo *repositories.ActivityRepository
}

func NewInterviewHandler(service *services.InterviewService, activityRepo *repositories.ActivityRepository) *InterviewHandler {
	return &InterviewHandler{service: service, activityRepo: activityRepo}
}

func (h *InterviewHandler) CreateReal(c *gin.Context) {
	var req struct {
		URL      string    `json:"url" binding:"required"`
		Company  string    `json:"company"`
		Position string    `json:"position"`
		Grade    string    `json:"grade"`
		Stack    string    `json:"stack"`
		Date     time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	studentID := c.GetString("userID")
	if studentID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	input := services.InterviewInput{
		URL:      req.URL,
		Company:  req.Company,
		Position: req.Position,
		Grade:    req.Grade,
		Stack:    req.Stack,
		Date:     req.Date,
	}
	interview, err := h.service.CreateReal(studentID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, interview)
}

func (h *InterviewHandler) CreateMock(c *gin.Context) {
	var req struct {
		StudentID string    `json:"student_id" binding:"required"`
		URL       string    `json:"url"`
		Company   string    `json:"company"`
		Position  string    `json:"position"`
		Grade     string    `json:"grade"`
		Stack     string    `json:"stack"`
		Date      time.Time `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	buddyID := c.GetString("userID")
	input := services.InterviewInput{
		URL:      req.URL,
		Company:  req.Company,
		Position: req.Position,
		Grade:    req.Grade,
		Stack:    req.Stack,
		Date:     req.Date,
	}
	interview, err := h.service.CreateMock(buddyID, req.StudentID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	metadata := `{"company": "` + req.Company + `", "position": "` + req.Position + `"}`
	_ = h.activityRepo.CreateActivity(req.StudentID, buddyID, "mock_created", "interview", interview.ID, metadata)

	c.JSON(http.StatusCreated, interview)
}

func (h *InterviewHandler) CompleteMock(c *gin.Context) {
	var req struct {
		InterviewID string `json:"interview_id" binding:"required"`
		Feedback    string `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CompleteMock(req.InterviewID, req.Feedback); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "completed"})
}

func (h *InterviewHandler) MyInterviews(c *gin.Context) {
	userID := c.GetString("userID")
	studentInterviews, _ := h.service.GetUserInterviews(userID, false)
	buddyInterviews, _ := h.service.GetUserInterviews(userID, true)
	all := append(studentInterviews, buddyInterviews...)
	c.JSON(http.StatusOK, all)
}

func (h *InterviewHandler) AllReal(c *gin.Context) {
	interviews, err := h.service.GetAllRealInterviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, interviews)
}
