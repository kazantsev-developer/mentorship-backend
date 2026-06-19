package models

import (
	"time"
)

// MaterialType defines the type of a learning material
type MaterialType string

const (
	MaterialTheory    MaterialType = "theory"
	MaterialQuestions MaterialType = "questions"
	MaterialPractice  MaterialType = "practice"
	MaterialHomework  MaterialType = "homework"
)

// Material represents a learning material within a roadmap block
type Material struct {
	ID           string       `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BlockID      string       `json:"block_id" gorm:"type:uuid;not null;index"`
	Title        string       `json:"title" gorm:"not null"`
	Description  string       `json:"description"`
	Type         MaterialType `json:"type" gorm:"type:varchar(20);not null"`
	ContentType  string       `json:"content_type" gorm:"type:varchar(50)"`
	URL          string       `json:"url"`
	Content      string       `json:"content" gorm:"type:text"`
	PreviewTitle string       `json:"preview_title"`
	PreviewDesc  string       `json:"preview_desc"`
	PreviewImage string       `json:"preview_image"`
	Source       string       `json:"source"`
	IsRequired   bool         `json:"is_required" gorm:"default:false"`
	IsActive     bool         `json:"is_active" gorm:"default:true"`
	SortOrder    int          `json:"sort_order" gorm:"default:0"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty" gorm:"index"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
