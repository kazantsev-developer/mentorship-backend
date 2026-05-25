package models

import (
	"time"
)

type Achievement struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string `gorm:"not null"`
	Description     string
	RewardBonus     int `gorm:"default:0"`
	ImageURL        string
	ConditionType   string `gorm:"type:varchar(50);not null"`
	ConditionParams string `gorm:"type:text"`
	IsActive        bool   `gorm:"default:true"`
	SortOrder       int    `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserAchievement struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID        string    `gorm:"type:uuid;not null;index"`
	AchievementID string    `gorm:"type:uuid;not null;index"`
	ReceivedAt    time.Time `gorm:"not null"`
	CreatedAt     time.Time
}
