// Package config provides application configuration loading from environment
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config aggregates all configuration sections
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Admin    AdminConfig
}

// ServerConfig contains server-related settings
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig contains database connection parameters
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig contains JSON Web Token settings
type JWTConfig struct {
	Secret      string
	ExpiresHour int
	ExpiresAt   time.Duration
}

// AdminConfig contains default admin credentials
type AdminConfig struct {
	Login    string
	Password string
}

// AppConfig holds the loaded application configuration
var AppConfig *Config

// Load reads configuration from environment variables and optional .env file
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("no .env file found, using system env")
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "mentor"),
			Password: getEnv("DB_PASSWORD", "secret"),
			DBName:   getEnv("DB_NAME", "mentorship"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "default-secret-change-me"),
			ExpiresHour: getEnvAsInt("JWT_EXPIRES_HOURS", 24),
		},
		Admin: AdminConfig{
			Login:    getEnv("ADMIN_LOGIN", "admin"),
			Password: getEnv("ADMIN_PASSWORD", "admin123"),
		},
	}

	cfg.JWT.ExpiresAt = time.Duration(cfg.JWT.ExpiresHour) * time.Hour

	AppConfig = cfg
	return cfg, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
