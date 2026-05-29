package models

import "time"

type FinalCheckType string

const (
	FinalTechnical FinalCheckType = "final_technical"
	FinalRoast     FinalCheckType = "final_roast"
)

type FinalCheck struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	StudentID   string         `gorm:"type:uuid;not null;index" json:"student_id"`
	BuddyID     string         `gorm:"type:uuid;not null;index" json:"buddy_id"`
	Type        FinalCheckType `gorm:"type:varchar(30);not null" json:"type"`
	Status      string         `gorm:"type:varchar(30);default:'not_available'" json:"status"`
	ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
