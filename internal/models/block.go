package models

import (
	"time"
)

// Block represents a roadmap block containing materials
type Block struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Materials []Material `json:"materials,omitempty" gorm:"foreignKey:BlockID"`
}
