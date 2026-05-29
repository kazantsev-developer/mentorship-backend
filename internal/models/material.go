package models

import "time"

type MaterialType string

const (
	MaterialTheory    MaterialType = "theory"
	MaterialQuestions MaterialType = "questions"
	MaterialPractice  MaterialType = "practice"
	MaterialHomework  MaterialType = "homework"
)

type Material struct {
	ID           string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BlockID      string       `gorm:"type:uuid;not null;index" json:"block_id"`
	Title        string       `gorm:"not null" json:"title"`
	Description  string       `json:"description"`
	Type         MaterialType `gorm:"type:varchar(20);not null" json:"type"`
	ContentType  string       `gorm:"type:varchar(50)" json:"content_type"`
	URL          string       `json:"url"`
	Content      string       `gorm:"type:text" json:"content"`
	PreviewTitle string       `json:"preview_title"`
	PreviewDesc  string       `json:"preview_desc"`
	PreviewImage string       `json:"preview_image"`
	Source       string       `json:"source"`
	IsRequired   bool         `gorm:"default:false" json:"is_required"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
	SortOrder    int          `gorm:"default:0" json:"sort_order"`
	DeletedAt    *time.Time   `gorm:"index" json:"-"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
