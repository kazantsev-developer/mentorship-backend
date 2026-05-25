package models

type Role string

const (
	RoleStudent Role = "student"
	RoleBuddy   Role = "buddy"
	RoleAdmin   Role = "admin"
)

type UserRole struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"type:uuid;not null"`
	Role   Role   `gorm:"type:varchar(20);not null"`
}
