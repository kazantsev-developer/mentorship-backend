package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type OneOnOneHandler struct {
	service *services.OneOnOneService
}

func NewOneOnOneHandler(service *services.OneOnOneService) *OneOnOneHandler {
	return &OneOnOneHandler{service: service}
}

func (h *OneOnOneHandler) CreateRequest(c *gin.Context) {
	studentID := c.GetString("userID")
	req, err := h.service.CreateRequest(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func (h *OneOnOneHandler) ListRequests(c *gin.Context) {
	requests, err := h.service.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *OneOnOneHandler) Approve(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminID := c.GetString("userID")
	if err := h.service.ApproveRequest(req.RequestID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}

func (h *OneOnOneHandler) Reject(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	adminID := c.GetString("userID")
	if err := h.service.RejectRequest(req.RequestID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "rejected"})
}
