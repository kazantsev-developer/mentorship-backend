package models

// Role defines a user role within the platform
type Role string

const (
	RoleStudent Role = "student"
	RoleBuddy   Role = "buddy"
	RoleAdmin   Role = "admin"
)

// UserRole associates a role with a user
type UserRole struct {
	ID     uint   `gorm:"primaryKey"`
	UserID string `gorm:"type:uuid;not null"`
	Role   Role   `gorm:"type:varchar(20);not null"`
}
