// Package db provides database connection and migration utilities
package db

import (
	"fmt"
	"log"

	"github.com/kazantsev/mentorship-backend/internal/config"
	"github.com/kazantsev/mentorship-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the global database instance
var DB *gorm.DB

// InitDB initializes the database connection and runs auto-migration
func InitDB(cfg *config.Config) error {
	dsn := cfg.GetDSN()
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := DB.AutoMigrate(
		&models.User{},
		&models.UserRole{},
		&models.Block{},
		&models.Material{},
		&models.MaterialProgress{},
		&models.BlockProgress{},
		&models.BonusBalance{},
		&models.BonusTransaction{},
		&models.Achievement{},
		&models.UserAchievement{},
		&models.CalendarEvent{},
		&models.FinalCheck{},
		&models.Interview{},
		&models.OneOnOneRequest{},
		&models.StudentBuddyAssignment{},
		&models.ActivityEvent{},
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	log.Println("database connected and migrated successfully")
	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return DB
}
