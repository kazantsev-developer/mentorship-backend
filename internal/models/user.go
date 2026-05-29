package models

import (
	"time"
)

type User struct {
	ID                string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Login             string     `gorm:"unique;not null" json:"login"`
	PasswordHash      string     `gorm:"not null" json:"-"`
	DisplayName       string     `gorm:"not null" json:"display_name"`
	AvatarURL         string     `json:"avatar_url"`
	About             string     `json:"about"`
	TelegramUsername  string     `gorm:"unique" json:"telegram_username"`
	LearningStartedAt *time.Time `json:"learning_started_at"`
	IsProfilePrivate  bool       `gorm:"default:false" json:"is_profile_private"`
	IsDeleted         bool       `gorm:"default:false" json:"-"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type BuddyStudentResponse struct {
	User
	CurrentBlockTitle string `json:"current_block_title"`
	ProgressPercent   int    `json:"progress_percent"`
	Status            string `json:"status"`
	DaysInactive      int    `json:"days_inactive"`
}
