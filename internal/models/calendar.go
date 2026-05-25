package models

import "time"

type CalendarEvent struct {
	ID              string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string `gorm:"not null"`
	StudentID       string `gorm:"type:uuid;not null;index"`
	BuddyID         string `gorm:"type:uuid;not null;index"`
	Type            string `gorm:"type:varchar(50)"`
	StartDatetime   time.Time
	EndDatetime     *time.Time
	Description     string
	ReminderEnabled bool `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
