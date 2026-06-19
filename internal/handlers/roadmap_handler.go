package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

// RoadmapHandler handles roadmap retrieval requests
type RoadmapHandler struct {
	roadmapService *services.RoadmapService
}

// NewRoadmapHandler returns a new RoadmapHandler instance
func NewRoadmapHandler(roadmapService *services.RoadmapService) *RoadmapHandler {
	return &RoadmapHandler{roadmapService: roadmapService}
}

// GetRoadmap returns the full learning roadmap for the current student
func (h *RoadmapHandler) GetRoadmap(c *gin.Context) {
	studentID := c.GetString("userID")
	roadmap, err := h.roadmapService.GetFullRoadmap(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch roadmap"})
		return
	}
	c.JSON(http.StatusOK, roadmap)
}
