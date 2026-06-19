package models

import (
	"time"
)

// BlockStatus defines the progress status of a roadmap block for a student
type BlockStatus string

const (
	BlockNotStarted               BlockStatus = "not_started"
	BlockInProgress               BlockStatus = "in_progress"
	BlockWaitingBuddyConfirmation BlockStatus = "waiting_buddy_confirmation"
	BlockApproved                 BlockStatus = "approved"
)

// MaterialProgress records a student's viewed material
type MaterialProgress struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string    `gorm:"type:uuid;not null;index"`
	MaterialID string    `gorm:"type:uuid;not null;index"`
	ViewedAt   time.Time `gorm:"not null"`
	CreatedAt  time.Time
}

// BlockProgress tracks a student's block status and approval
type BlockProgress struct {
	ID         string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string      `gorm:"type:uuid;not null;index"`
	BlockID    string      `gorm:"type:uuid;not null;index"`
	Status     BlockStatus `gorm:"type:varchar(30);default:'not_started'"`
	ApprovedBy *string     `gorm:"type:uuid"`
	ApprovedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
