package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string `yaml:"app_env"`
	PostgresHost     string `yaml:"postgres_host"`
	PostgresPort     int    `yaml:"postgres_port"`
	PostgresUser     string `yaml:"postgres_user"`
	PostgresPassword string `yaml:"postgres_password"`
	PostgresDBName   string `yaml:"postgres_dbname"`
	PostgresSSLMode  string `yaml:"postgres_sslmode"`
	PostgresMaxConns int    `yaml:"postgres_max_conns"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:           getEnv("APP_ENV", "dev"),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvAsInt("POSTGRES_PORT", 5432),
		PostgresUser:     getEnv("POSTGRES_USER", "oom_user_LQkR"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "r9VYUDHXxRk"),
		PostgresDBName:   getEnv("POSTGRES_DBNAME", "oom_db"),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
		PostgresMaxConns: getEnvAsInt("POSTGRES_MAX_CONNS", 20),
	}

	if cfg.PostgresHost == "" || cfg.PostgresPort == 0 || cfg.PostgresUser == "" || cfg.PostgresPassword == "" || cfg.PostgresDBName == "" || cfg.PostgresSSLMode == "" {
		return nil, fmt.Errorf("missing postgres configuration")
	}

	return cfg, nil
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
