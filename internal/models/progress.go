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
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string    `json:"student_id" gorm:"type:uuid;not null;index"`
	MaterialID string    `json:"material_id" gorm:"type:uuid;not null;index"`
	ViewedAt   time.Time `json:"viewed_at" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

// BlockProgress tracks a student's block status and approval
type BlockProgress struct {
	ID         string      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string      `json:"student_id" gorm:"type:uuid;not null;index"`
	BlockID    string      `json:"block_id" gorm:"type:uuid;not null;index"`
	Status     BlockStatus `json:"status" gorm:"type:varchar(30);default:'not_started'"`
	ApprovedBy *string     `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt *time.Time  `json:"approved_at"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
