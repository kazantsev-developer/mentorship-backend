package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"github.com/kazantsev/mentorship-backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type registerInput struct {
	Login            string   `json:"login" binding:"required"`
	Password         string   `json:"password" binding:"required"`
	DisplayName      string   `json:"display_name" binding:"required"`
	TelegramUsername string   `json:"telegram_username"`
	Roles            []string `json:"roles" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var roles []models.Role
	for _, r := range input.Roles {
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

	req := services.RegisterRequest{
		Login:            input.Login,
		Password:         input.Password,
		DisplayName:      input.DisplayName,
		TelegramUsername: input.TelegramUsername,
		Roles:            roles,
	}

	user, err := h.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type loginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authService.Login(services.LoginRequest{
		Login:    input.Login,
		Password: input.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
		"user":  resp.User,
		"roles": resp.Roles,
	})
}
