package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

// CalendarHandler handles calendar event operations
type CalendarHandler struct {
	service *services.CalendarService
}

// NewCalendarHandler returns a new CalendarHandler instance
func NewCalendarHandler(service *services.CalendarService) *CalendarHandler {
	return &CalendarHandler{service: service}
}

// CreateEvent creates a new calendar event
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	buddyID := c.GetString("userID")
	event, err := h.service.CreateEvent(buddyID, req.StudentID, req.Title, req.Type, req.Description, &req.StartDatetime, req.EndDatetime, req.ReminderEnabled)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}
	c.JSON(http.StatusCreated, event)
}

// MyEvents returns events for the current user (student or buddy)
func (h *CalendarHandler) MyEvents(c *gin.Context) {
	userID := c.GetString("userID")
	events, err := h.service.GetEventsForStudent(userID)
	if err != nil {
		events, err = h.service.GetEventsForBuddy(userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

// UpcomingEvents returns events within the upcoming N days (default 7)
func (h *CalendarHandler) UpcomingEvents(c *gin.Context) {
	userID := c.GetString("userID")
	days := 7
	events, err := h.service.GetUpcomingEvents(userID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch upcoming events"})
		return
	}
	c.JSON(http.StatusOK, events)
}
