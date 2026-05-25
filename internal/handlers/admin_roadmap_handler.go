package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type AdminRoadmapHandler struct {
	db *gorm.DB
}

func NewAdminRoadmapHandler(db *gorm.DB) *AdminRoadmapHandler {
	return &AdminRoadmapHandler{db: db}
}

func (h *AdminRoadmapHandler) ListBlocks(c *gin.Context) {
	var blocks []models.Block
	h.db.Where("deleted_at IS NULL").Order("sort_order").Find(&blocks)
	c.JSON(http.StatusOK, blocks)
}

func (h *AdminRoadmapHandler) CreateBlock(c *gin.Context) {
	var block models.Block
	if err := c.ShouldBindJSON(&block); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	block.ID = uuid.New().String()
	block.CreatedAt = time.Now()
	block.UpdatedAt = time.Now()
	if err := h.db.Create(&block).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, block)
}

func (h *AdminRoadmapHandler) UpdateBlock(c *gin.Context) {
	id := c.Param("id")
	var block models.Block
	if err := h.db.First(&block, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "block not found"})
		return
	}
	var updates models.Block
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	block.Title = updates.Title
	block.Description = updates.Description
	block.SortOrder = updates.SortOrder
	block.IsActive = updates.IsActive
	block.UpdatedAt = time.Now()
	if err := h.db.Save(&block).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, block)
}

func (h *AdminRoadmapHandler) DeleteBlock(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	if err := h.db.Model(&models.Block{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
