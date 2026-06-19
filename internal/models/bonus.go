package models

import (
	"time"
)

// BonusTransactionType defines the type of a bonus transaction
type BonusTransactionType string

const (
	BonusTypeAchievement   BonusTransactionType = "achievement_reward"
	BonusTypeDiscount      BonusTransactionType = "discount_conversion"
	BonusTypeOneOnOneSpend BonusTransactionType = "one_on_one_spend"
	BonusTypeManual        BonusTransactionType = "manual_adjustment"
	BonusTypeRefund        BonusTransactionType = "refund"
)

// BonusBalance stores the current bonus balance for a user
type BonusBalance struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"`
	Balance   int       `json:"balance" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BonusTransaction records an individual bonus transaction
type BonusTransaction struct {
	ID         string               `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string               `json:"user_id" gorm:"type:uuid;not null;index"`
	Type       BonusTransactionType `json:"type" gorm:"type:varchar(30);not null"`
	Amount     int                  `json:"amount" gorm:"not null"`
	Reason     string               `json:"reason" gorm:"type:text"`
	SourceType string               `json:"source_type" gorm:"type:varchar(50)"`
	SourceID   string               `json:"source_id" gorm:"type:uuid"`
	CreatedAt  time.Time            `json:"created_at"`
}
