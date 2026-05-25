package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type CalendarHandler struct {
	service *services.CalendarService
}

func NewCalendarHandler(service *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

func (h *CalendarHandler) CreateEvent(c *gin.Context) {
	var req struct {
		StudentID       string     `json:"student_id" binding:"required"`
		Title           string     `json:"title" binding:"required"`
		Type            string     `json:"type"`
		Description     string     `json:"description"`
		StartDatetime   time.Time  `json:"start_datetime" binding:"required"`
		EndDatetime     *time.Time `json:"end_datetime"`
		ReminderEnabled bool       `json:"reminder_enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	buddyID := c.GetString("userID")
	event, err := h.service.CreateEvent(buddyID, req.StudentID, req.Title, req.Type, req.Description, &req.StartDatetime, req.EndDatetime, req.ReminderEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func (h *CalendarHandler) MyEvents(c *gin.Context) {
	userID := c.GetString("userID")
	events, err := h.service.GetEventsForStudent(userID)
	if err != nil {
		events, err = h.service.GetEventsForBuddy(userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *CalendarHandler) UpcomingEvents(c *gin.Context) {
	userID := c.GetString("userID")
	days := 7
	if d, ok := c.GetQuery("days"); ok {
		_ = d
	}
	events, err := h.service.GetUpcomingEvents(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}
