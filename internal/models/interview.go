package models

import "time"

// InterviewType defines the type of an interview
type InterviewType string

const (
	InterviewMock InterviewType = "mock"
	InterviewReal InterviewType = "real"
)

// Interview represents an interview record (mock or real)
type Interview struct {
	ID        string        `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Type      InterviewType `json:"type" gorm:"type:varchar(10);not null"`
	StudentID string        `json:"student_id" gorm:"type:uuid;not null;index"`
	BuddyID   *string       `json:"buddy_id" gorm:"type:uuid;index"`
	URL       string        `json:"url"`
	Company   string        `json:"company"`
	Position  string        `json:"position"`
	Grade     string        `json:"grade"`
	Stack     string        `json:"stack"`
	Date      time.Time     `json:"date"`
	Status    string        `json:"status" gorm:"type:varchar(30)"`
	Result    string        `json:"result" gorm:"type:varchar(30)"`
	Feedback  string        `json:"feedback" gorm:"type:text"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
