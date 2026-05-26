package models

import "time"

type InterviewType string

const (
	InterviewMock InterviewType = "mock"
	InterviewReal InterviewType = "real"
)

type Interview struct {
	ID        string        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Type      InterviewType `gorm:"type:varchar(10);not null"`
	StudentID string        `gorm:"type:uuid;not null;index"`
	BuddyID   *string       `gorm:"type:uuid;index"` // разрешаем NULL
	URL       string
	Company   string
	Position  string
	Grade     string
	Stack     string
	Date      time.Time
	Status    string `gorm:"type:varchar(30)"`
	Result    string `gorm:"type:varchar(30)"`
	Feedback  string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
