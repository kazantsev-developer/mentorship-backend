package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/gorm"
)

type AdminMaterialHandler struct {
	db *gorm.DB
}

func NewAdminMaterialHandler(db *gorm.DB) *AdminMaterialHandler {
	return &AdminMaterialHandler{db: db}
}

func (h *AdminMaterialHandler) ListMaterials(c *gin.Context) {
	var materials []models.Material
	blockID := c.Query("block_id")
	query := h.db.Where("deleted_at IS NULL")
	if blockID != "" {
		query = query.Where("block_id = ?", blockID)
	}
	query.Order("sort_order").Find(&materials)
	c.JSON(http.StatusOK, materials)
}

func (h *AdminMaterialHandler) CreateMaterial(c *gin.Context) {
	var material models.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	material.ID = uuid.New().String()
	material.CreatedAt = time.Now()
	material.UpdatedAt = time.Now()
	if err := h.db.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, material)
}

func (h *AdminMaterialHandler) UpdateMaterial(c *gin.Context) {
	id := c.Param("id")
	var material models.Material
	if err := h.db.First(&material, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "material not found"})
		return
	}
	var updates models.Material
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	material.Title = updates.Title
	material.Description = updates.Description
	material.Type = updates.Type
	material.ContentType = updates.ContentType
	material.URL = updates.URL
	material.Content = updates.Content
	material.IsRequired = updates.IsRequired
	material.IsActive = updates.IsActive
	material.SortOrder = updates.SortOrder
	material.UpdatedAt = time.Now()
	if err := h.db.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, material)
}

func (h *AdminMaterialHandler) DeleteMaterial(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	if err := h.db.Model(&models.Material{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminMaterialHandler) ToggleMaterialStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.db.Model(&models.Material{}).Where("id = ?", id).Update("is_active", req.IsActive).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
