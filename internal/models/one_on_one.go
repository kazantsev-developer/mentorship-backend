package models

import "time"

// OneOnOneRequest represents a student's request for a one-on-one meeting
type OneOnOneRequest struct {
	ID          string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID   string     `json:"student_id" gorm:"type:uuid;not null;index"`
	Status      string     `json:"status" gorm:"type:varchar(30);default:'pending'"`
	ProcessedBy *string    `json:"processed_by" gorm:"type:uuid"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
