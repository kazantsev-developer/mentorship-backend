package models

import "time"

type OneOnOneRequest struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	StudentID   string    `gorm:"type:uuid;not null;index"`
	Status      string    `gorm:"type:varchar(30);default:'pending'"`
	ProcessedBy *string   `gorm:"type:uuid"`
	ProcessedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
