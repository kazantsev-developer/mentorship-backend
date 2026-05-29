package models

import "time"

type CalendarEvent struct {
	ID              string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title           string     `gorm:"not null" json:"title"`
	StudentID       string     `gorm:"type:uuid;not null;index" json:"student_id"`
	BuddyID         string     `gorm:"type:uuid;not null;index" json:"buddy_id"`
	Type            string     `gorm:"type:varchar(50)" json:"type"`
	StartDatetime   time.Time  `json:"start_datetime"`
	EndDatetime     *time.Time `json:"end_datetime,omitempty"`
	Description     string     `json:"description"`
	ReminderEnabled bool       `gorm:"default:false" json:"reminder_enabled"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
