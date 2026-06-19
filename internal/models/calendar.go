package models

import "time"

// CalendarEvent represents a scheduled calendar event for a student and buddy
type CalendarEvent struct {
	ID              string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title           string     `json:"title" gorm:"not null"`
	StudentID       string     `json:"student_id" gorm:"type:uuid;not null;index"`
	BuddyID         string     `json:"buddy_id" gorm:"type:uuid;not null;index"`
	Type            string     `json:"type" gorm:"type:varchar(50)"`
	StartDatetime   time.Time  `json:"start_datetime"`
	EndDatetime     *time.Time `json:"end_datetime"`
	Description     string     `json:"description"`
	ReminderEnabled bool       `json:"reminder_enabled" gorm:"default:false"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
