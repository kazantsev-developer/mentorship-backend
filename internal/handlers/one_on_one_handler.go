package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

// OneOnOneHandler handles one-on-one meeting requests
type OneOnOneHandler struct {
	service *services.OneOnOneService
}

// NewOneOnOneHandler returns a new OneOnOneHandler instance
func NewOneOnOneHandler(service *services.OneOnOneService) *OneOnOneHandler {
	return &OneOnOneHandler{service: service}
}

// CreateRequest creates a new one-on-one meeting request for the authenticated student
func (h *OneOnOneHandler) CreateRequest(c *gin.Context) {
	studentID := c.GetString("userID")
	req, err := h.service.CreateRequest(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// ListRequests returns all one-on-one requests
func (h *OneOnOneHandler) ListRequests(c *gin.Context) {
	requests, err := h.service.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch requests"})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// Approve approves a pending one-on-one request
func (h *OneOnOneHandler) Approve(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	adminID := c.GetString("userID")
	if err := h.service.ApproveRequest(req.RequestID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to approve request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

// Reject rejects a pending one-on-one request
func (h *OneOnOneHandler) Reject(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	adminID := c.GetString("userID")
	if err := h.service.RejectRequest(req.RequestID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to reject request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}
