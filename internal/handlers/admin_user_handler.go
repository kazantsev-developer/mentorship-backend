package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"github.com/kazantsev/mentorship-backend/internal/utils"
	"gorm.io/gorm"
)

type AdminUserHandler struct {
	userRepo     *repositories.UserRepository
	authService  *services.AuthService
	progressRepo *repositories.ProgressRepository
	blockRepo    *repositories.BlockRepository
	activityRepo *repositories.ActivityRepository
	db           *gorm.DB
}

func NewAdminUserHandler(
	userRepo *repositories.UserRepository,
	authService *services.AuthService,
	progressRepo *repositories.ProgressRepository,
	blockRepo *repositories.BlockRepository,
	activityRepo *repositories.ActivityRepository,
	db *gorm.DB,
) *AdminUserHandler {
	return &AdminUserHandler{
		userRepo:     userRepo,
		authService:  authService,
		progressRepo: progressRepo,
		blockRepo:    blockRepo,
		activityRepo: activityRepo,
		db:           db,
	}
}

// ListUsers возвращает список пользователей с их ролями
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	includeDeleted := c.Query("deleted") == "true"
	query := h.userRepo.GetDB().Model(&models.User{})
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Структура для ответа с ролями
	type UserWithRoles struct {
		models.User
		Roles []string `json:"roles"`
	}
	result := make([]UserWithRoles, 0, len(users))

	for _, u := range users {
		var roles []string
		// Загружаем роли пользователя из таблицы user_roles
		h.db.Table("user_roles").Where("user_id = ?", u.ID).Pluck("role", &roles)
		u.PasswordHash = "" // не показываем хэш пароля
		result = append(result, UserWithRoles{
			User:  u,
			Roles: roles,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminUserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Login             string   `json:"login" binding:"required"`
		Password          string   `json:"password" binding:"required"`
		DisplayName       string   `json:"display_name" binding:"required"`
		TelegramUsername  string   `json:"telegram_username"`
		Roles             []string `json:"roles" binding:"required"`
		LearningStartedAt *string  `json:"learning_started_at"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var roles []models.Role
	for _, r := range req.Roles {
		switch r {
		case "student":
			roles = append(roles, models.RoleStudent)
		case "buddy":
			roles = append(roles, models.RoleBuddy)
		case "admin":
			roles = append(roles, models.RoleAdmin)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
	}
	regReq := services.RegisterRequest{
		Login:            req.Login,
		Password:         req.Password,
		DisplayName:      req.DisplayName,
		TelegramUsername: req.TelegramUsername,
		Roles:            roles,
	}
	user, err := h.authService.Register(regReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.LearningStartedAt != nil && *req.LearningStartedAt != "" {
		t, parseErr := time.Parse(time.RFC3339, *req.LearningStartedAt)
		if parseErr == nil {
			user.LearningStartedAt = &t
			_ = h.userRepo.Update(user)
		}
	}
	user.PasswordHash = ""
	c.JSON(http.StatusCreated, user)
}

func (h *AdminUserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("user_id")
	var updates struct {
		DisplayName       string  `json:"display_name"`
		AvatarURL         string  `json:"avatar_url"`
		About             string  `json:"about"`
		TelegramUsername  string  `json:"telegram_username"`
		IsProfilePrivate  *bool   `json:"is_profile_private"`
		LearningStartedAt *string `json:"learning_started_at"`
	}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if updates.DisplayName != "" {
		user.DisplayName = updates.DisplayName
	}
	if updates.AvatarURL != "" {
		user.AvatarURL = updates.AvatarURL
	}
	if updates.About != "" {
		user.About = updates.About
	}
	if updates.TelegramUsername != "" {
		user.TelegramUsername = updates.TelegramUsername
	}
	if updates.IsProfilePrivate != nil {
		user.IsProfilePrivate = *updates.IsProfilePrivate
	}
	if updates.LearningStartedAt != nil && *updates.LearningStartedAt != "" {
		t, parseErr := time.Parse(time.RFC3339, *updates.LearningStartedAt)
		if parseErr == nil {
			user.LearningStartedAt = &t
		}
	}
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	if err := h.userRepo.SoftDelete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *AdminUserHandler) ChangePassword(c *gin.Context) {
	userID := c.Param("user_id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = hashed
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}

func (h *AdminUserHandler) GetUser(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *AdminUserHandler) GetUserProgress(c *gin.Context) {
	userID := c.Param("user_id")
	blocks, err := h.blockRepo.GetAllActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type BlockProgress struct {
		BlockID string `json:"block_id"`
		Title   string `json:"block_title"`
		Status  string `json:"status"`
		Percent int    `json:"percent"`
	}
	result := make([]BlockProgress, 0, len(blocks))
	for _, block := range blocks {
		status, _ := h.progressRepo.GetBlockStatus(userID, block.ID)
		percent, _ := h.progressRepo.GetBlockProgressPercent(userID, block.ID)
		result = append(result, BlockProgress{
			BlockID: block.ID,
			Title:   block.Title,
			Status:  status,
			Percent: percent,
		})
	}
	c.JSON(http.StatusOK, result)
}

func (h *AdminUserHandler) ApproveBlock(c *gin.Context) {
	userID := c.Param("user_id")
	blockID := c.Param("block_id")
	adminID := c.GetString("userID")
	if err := h.progressRepo.ConfirmBlock(userID, blockID, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	metadata := `{"block_id": "` + blockID + `", "approved_by_admin": true}`
	_ = h.activityRepo.CreateActivity(userID, adminID, "block_approved", "block", blockID, metadata)
	c.JSON(http.StatusOK, gin.H{"status": "approved"})
}
