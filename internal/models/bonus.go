package models

import (
	"time"
)

type BonusTransactionType string

const (
	BonusTypeAchievement   BonusTransactionType = "achievement_reward"
	BonusTypeDiscount      BonusTransactionType = "discount_conversion"
	BonusTypeOneOnOneSpend BonusTransactionType = "one_on_one_spend"
	BonusTypeManual        BonusTransactionType = "manual_adjustment"
	BonusTypeRefund        BonusTransactionType = "refund"
)

type BonusBalance struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex"`
	Balance   int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BonusTransaction struct {
	ID         string               `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string               `gorm:"type:uuid;not null;index"`
	Type       BonusTransactionType `gorm:"type:varchar(30);not null"`
	Amount     int                  `gorm:"not null"`
	Reason     string               `gorm:"type:text"`
	SourceType string               `gorm:"type:varchar(50)"`
	SourceID   string               `gorm:"type:uuid"`
	CreatedAt  time.Time
}
