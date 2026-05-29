package models

import "time"

type ActivityEvent struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID     string    `gorm:"type:uuid;not null;index" json:"user_id"`
	ActorID    *string   `gorm:"type:uuid" json:"actor_id,omitempty"`
	Type       string    `gorm:"type:varchar(50);not null" json:"type"`
	SourceType string    `gorm:"type:varchar(50)" json:"source_type"`
	SourceID   string    `gorm:"type:uuid" json:"source_id"`
	Metadata   string    `gorm:"type:text" json:"metadata"`
	CreatedAt  time.Time `json:"created_at"`
}
