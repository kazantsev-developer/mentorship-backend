package models

import (
	"time"
)

// Block represents a roadmap block containing materials
type Block struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string `gorm:"not null"`
	Description string
	SortOrder   int        `gorm:"default:0"`
	IsActive    bool       `gorm:"default:true"`
	DeletedAt   *time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Materials []Material `gorm:"foreignKey:BlockID"`
}
