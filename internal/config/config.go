package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	Admin          AdminConfig
	AllowedOrigins []string
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret      string
	ExpiresHour int
	ExpiresAt   time.Duration
}

type AdminConfig struct {
	Login    string
	Password string
}

var AppConfig *Config

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system env")
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
		AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000")),
	}
	cfg.JWT.ExpiresAt = time.Duration(cfg.JWT.ExpiresHour) * time.Hour

	AppConfig = cfg
	return cfg, nil
}

func (c *Config) GetDSN() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}

func parseOrigins(originsStr string) []string {
	if originsStr == "" {
		return []string{}
	}
	parts := strings.Split(originsStr, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
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
