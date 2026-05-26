package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type RoadmapHandler struct {
	roadmapService *services.RoadmapService
}

func NewRoadmapHandler(roadmapService *services.RoadmapService) *RoadmapHandler {
	return &RoadmapHandler{roadmapService: roadmapService}
}

func (h *RoadmapHandler) GetRoadmap(c *gin.Context) {
	studentID := c.GetString("userID")
	roadmap, err := h.roadmapService.GetFullRoadmap(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch roadmap"})
		return
	}
	c.JSON(http.StatusOK, roadmap)
}
