package models

import "time"

type Block struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	SortOrder   int        `gorm:"default:0" json:"sort_order"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	DeletedAt   *time.Time `gorm:"index" json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Materials []Material `gorm:"foreignKey:BlockID" json:"materials,omitempty"`
}
