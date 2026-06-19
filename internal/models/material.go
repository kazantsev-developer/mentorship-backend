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
	ID           string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	BlockID      string `gorm:"type:uuid;not null;index"`
	Title        string `gorm:"not null"`
	Description  string
	Type         MaterialType `gorm:"type:varchar(20);not null"`
	ContentType  string       `gorm:"type:varchar(50)"`
	URL          string
	Content      string `gorm:"type:text"`
	PreviewTitle string
	PreviewDesc  string
	PreviewImage string
	Source       string
	IsRequired   bool       `gorm:"default:false"`
	IsActive     bool       `gorm:"default:true"`
	SortOrder    int        `gorm:"default:0"`
	DeletedAt    *time.Time `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
