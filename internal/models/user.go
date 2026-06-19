package models

import (
	"time"
)

// User represents a registered user of the mentorship platform
type User struct {
	ID                string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Login             string     `json:"login" gorm:"unique;not null"`
	PasswordHash      string     `json:"-" gorm:"not null"`
	DisplayName       string     `json:"display_name" gorm:"not null"`
	AvatarURL         string     `json:"avatar_url"`
	About             string     `json:"about"`
	TelegramUsername  string     `json:"telegram_username" gorm:"unique"`
	LearningStartedAt *time.Time `json:"learning_started_at"`
	IsProfilePrivate  bool       `json:"is_profile_private" gorm:"default:false"`
	IsDeleted         bool       `json:"is_deleted" gorm:"default:false"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
