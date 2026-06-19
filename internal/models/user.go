package models

import (
	"time"
)

// User represents a registered user of the mentorship platform
type User struct {
	ID                string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Login             string `gorm:"unique;not null"`
	PasswordHash      string `gorm:"not null"`
	DisplayName       string `gorm:"not null"`
	AvatarURL         string
	About             string
	TelegramUsername  string     `gorm:"unique"`
	LearningStartedAt *time.Time `json:"learning_started_at"`
	IsProfilePrivate  bool       `gorm:"default:false"`
	IsDeleted         bool       `gorm:"default:false"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
