package models

import "time"

type ActivityEvent struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string  `gorm:"type:uuid;not null;index"`
	ActorID    *string `gorm:"type:uuid"`
	Type       string  `gorm:"type:varchar(50);not null"`
	SourceType string  `gorm:"type:varchar(50)"`
	SourceID   string  `gorm:"type:uuid"`
	Metadata   string  `gorm:"type:text"`
	CreatedAt  time.Time
}
