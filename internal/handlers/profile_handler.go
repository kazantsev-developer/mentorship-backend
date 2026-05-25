package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kazantsev/mentorship-backend/internal/repositories"
)

type ProfileHandler struct {
	userRepo *repositories.UserRepository
}

func NewProfileHandler(userRepo *repositories.UserRepository) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	var updates struct {
		DisplayName      string `json:"display_name"`
		AvatarURL        string `json:"avatar_url"`
		About            string `json:"about"`
		IsProfilePrivate bool   `json:"is_profile_private"`
	}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if updates.DisplayName != "" {
		user.DisplayName = updates.DisplayName
	}
	user.AvatarURL = updates.AvatarURL
	user.About = updates.About
	user.IsProfilePrivate = updates.IsProfilePrivate
	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}

func (h *ProfileHandler) PublicProfile(c *gin.Context) {
	userID := c.Param("user_id")
	user, err := h.userRepo.FindByID(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.IsProfilePrivate {
		c.JSON(http.StatusOK, gin.H{
			"id":                user.ID,
			"display_name":      user.DisplayName,
			"avatar_url":        user.AvatarURL,
			"telegram_username": user.TelegramUsername,
		})
		return
	}
	user.PasswordHash = ""
	c.JSON(http.StatusOK, user)
}
