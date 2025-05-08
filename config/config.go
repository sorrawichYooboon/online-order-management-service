package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv           string `yaml:"app_env"`
	DatabaseHost     string `yaml:"database_host"`
	DatabasePort     int    `yaml:"database_port"`
	DatabaseUser     string `yaml:"database_user"`
	DatabasePassword string `yaml:"database_password"`
	DatabaseDBName   string `yaml:"database_dbname"`
	DatabaseSSLMode  string `yaml:"database_sslmode"`
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	dbPort, _ := strconv.Atoi(os.Getenv("DATABASE_PORT"))

	cfg := &Config{
		AppEnv:           os.Getenv("APP_ENV"),
		DatabaseHost:     os.Getenv("DATABASE_HOST"),
		DatabasePort:     dbPort,
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseDBName:   os.Getenv("DATABASE_DBNAME"),
		DatabaseSSLMode:  os.Getenv("DATABASE_SSLMODE"),
	}

	if cfg.DatabaseHost == "" || cfg.DatabasePort == 0 || cfg.DatabaseUser == "" || cfg.DatabasePassword == "" || cfg.DatabaseDBName == "" || cfg.DatabaseSSLMode == "" {
		return nil, fmt.Errorf("missing database configuration")
	}

	return cfg, nil
}
