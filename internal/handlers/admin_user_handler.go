package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
	"github.com/kazantsev/mentorship-backend/internal/services"
	"github.com/kazantsev/mentorship-backend/internal/utils"
)

// AdminUserHandler handles administrative operations on user accounts
type AdminUserHandler struct {
	userRepo    *repositories.UserRepository
	authService *services.AuthService
}

// NewAdminUserHandler returns a new AdminUserHandler instance
func NewAdminUserHandler(userRepo *repositories.UserRepository, authService *services.AuthService) *AdminUserHandler {
	return &AdminUserHandler{userRepo: userRepo, authService: authService}
}

// ListUsers returns all users, optionally including deleted ones
func (h *AdminUserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	includeDeleted := c.Query("deleted") == "true"
	query := h.userRepo.GetDB().Model(&models.User{})
	if !includeDeleted {
		query = query.Where("is_deleted = ?", false)
	}
	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	for i := range users {
		users[i].PasswordHash = ""
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser registers a new user with the given profile and roles
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusCreated, user)
}

// UpdateUser modifies profile fields of an existing user
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

// DeleteUser soft-deletes a user by ID
func (h *AdminUserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("user_id")
	if err := h.userRepo.SoftDelete(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// ChangePassword sets a new password for the specified user
func (h *AdminUserHandler) ChangePassword(c *gin.Context) {
	userID := c.Param("user_id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.PasswordHash = hashed
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "password updated"})
}
