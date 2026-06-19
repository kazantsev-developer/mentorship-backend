// Package models defines data structures for the mentorship platform
package models

import (
	"time"
)

// Achievement represents an unlockable achievement with conditions and rewards
type Achievement struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string    `json:"title" gorm:"not null"`
	Description     string    `json:"description"`
	RewardBonus     int       `json:"reward_bonus" gorm:"default:0"`
	ImageURL        string    `json:"image_url"`
	ConditionType   string    `json:"condition_type" gorm:"type:varchar(50);not null"`
	ConditionParams string    `json:"condition_params" gorm:"type:text"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	SortOrder       int       `json:"sort_order" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserAchievement tracks which achievements have been awarded to a user
type UserAchievement struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        string    `json:"user_id" gorm:"type:uuid;not null;index"`
	AchievementID string    `json:"achievement_id" gorm:"type:uuid;not null;index"`
	ReceivedAt    time.Time `json:"received_at" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
}
