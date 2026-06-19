package models

import "time"

// FinalCheckType represents the type of a final check
type FinalCheckType string

const (
	FinalTechnical FinalCheckType = "final_technical"
	FinalRoast     FinalCheckType = "final_roast"
)

// FinalCheck stores information about a final assessment
type FinalCheck struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID   string         `json:"student_id" gorm:"type:uuid;not null;index"`
	BuddyID     string         `json:"buddy_id" gorm:"type:uuid;not null;index"`
	Type        FinalCheckType `json:"type" gorm:"type:varchar(30);not null"`
	Status      string         `json:"status" gorm:"type:varchar(30);default:'not_available'"`
	ScheduledAt *time.Time     `json:"scheduled_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
