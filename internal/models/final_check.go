package models

import "time"

type FinalCheckType string

const (
	FinalTechnical FinalCheckType = "final_technical"
	FinalRoast     FinalCheckType = "final_roast"
)

type FinalCheck struct {
	ID          string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID   string         `gorm:"type:uuid;not null;index"`
	BuddyID     string         `gorm:"type:uuid;not null;index"`
	Type        FinalCheckType `gorm:"type:varchar(30);not null"`
	Status      string         `gorm:"type:varchar(30);default:'not_available'"`
	ScheduledAt *time.Time
	CompletedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
