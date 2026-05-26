package models

import (
	"time"
)

type BlockStatus string

const (
	BlockNotStarted               BlockStatus = "not_started"
	BlockInProgress               BlockStatus = "in_progress"
	BlockWaitingBuddyConfirmation BlockStatus = "waiting_buddy_confirmation"
	BlockApproved                 BlockStatus = "approved"
)

// MaterialProgress – отметка о просмотре материала студентом
type MaterialProgress struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string    `gorm:"type:uuid;not null;index"`
	MaterialID string    `gorm:"type:uuid;not null;index"`
	ViewedAt   time.Time `gorm:"not null"`
	CreatedAt  time.Time
}

// BlockProgress – статус блока для студента
type BlockProgress struct {
	ID         string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID  string      `gorm:"type:uuid;not null;index"`
	BlockID    string      `gorm:"type:uuid;not null;index"`
	Status     BlockStatus `gorm:"type:varchar(30);default:'not_started'"`
	ApprovedBy *string     `gorm:"type:uuid"` // ID пользователя (buddy/admin), кто подтвердил
	ApprovedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
