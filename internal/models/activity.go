package models

import "time"

// ActivityEvent records an activity log entry for a user
type ActivityEvent struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id" gorm:"type:uuid;not null;index"`
	ActorID    *string   `json:"actor_id" gorm:"type:uuid"`
	Type       string    `json:"type" gorm:"type:varchar(50);not null"`
	SourceType string    `json:"source_type" gorm:"type:varchar(50)"`
	SourceID   string    `json:"source_id" gorm:"type:uuid"`
	Metadata   string    `json:"metadata" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
}
