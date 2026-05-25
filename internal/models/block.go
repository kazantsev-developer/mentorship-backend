package models

import (
	"time"
)

type Block struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string `gorm:"not null"`
	Description string
	SortOrder   int        `gorm:"default:0"`
	IsActive    bool       `gorm:"default:true"`
	DeletedAt   *time.Time `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Связь с материалами
	Materials []Material `gorm:"foreignKey:BlockID"`
}
