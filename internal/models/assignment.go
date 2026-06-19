package models

import "time"

// StudentBuddyAssignment links a student to a buddy
type StudentBuddyAssignment struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID string `gorm:"type:uuid;not null;index"`
	BuddyID   string `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
}
