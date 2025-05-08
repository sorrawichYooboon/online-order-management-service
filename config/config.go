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
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	dbPort, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))

	cfg := &Config{
		AppEnv:           os.Getenv("APP_ENV"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     dbPort,
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDBName:   os.Getenv("POSTGRES_DBNAME"),
		PostgresSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

	if cfg.PostgresHost == "" || cfg.PostgresPort == 0 || cfg.PostgresUser == "" || cfg.PostgresPassword == "" || cfg.PostgresDBName == "" || cfg.PostgresSSLMode == "" {
		return nil, fmt.Errorf("missing postgres configuration")
	}

	return cfg, nil
}
