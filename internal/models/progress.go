package models

import (
	"time"
)

type BlockStatus string

const (
	BlockStatusNotStarted               BlockStatus = "not_started"
	BlockStatusInProgress               BlockStatus = "in_progress"
	BlockStatusWaitingBuddyConfirmation BlockStatus = "waiting_buddy_confirmation"
	BlockStatusApproved                 BlockStatus = "approved"
)

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

type MaterialProgress struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string `gorm:"type:uuid;not null;index"`
	MaterialID string `gorm:"type:uuid;not null;index"`
	ViewedAt   *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
