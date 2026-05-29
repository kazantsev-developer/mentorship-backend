package models

import "time"

type Achievement struct {
	ID              string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title           string    `gorm:"not null" json:"title"`
	Description     string    `json:"description"`
	RewardBonus     int       `gorm:"default:0" json:"reward_bonus"`
	ImageURL        string    `json:"image_url"`
	ConditionType   string    `gorm:"type:varchar(50);not null" json:"condition_type"`
	ConditionParams string    `gorm:"type:text" json:"condition_params"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	SortOrder       int       `gorm:"default:0" json:"sort_order"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserAchievement struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        string    `gorm:"type:uuid;not null;index" json:"user_id"`
	AchievementID string    `gorm:"type:uuid;not null;index" json:"achievement_id"`
	ReceivedAt    time.Time `gorm:"not null" json:"received_at"`
	CreatedAt     time.Time `json:"created_at"`
}
