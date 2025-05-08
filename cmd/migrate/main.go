package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sorrawichYooboon/online-order-management-service/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseDBName,
		cfg.DatabaseSSLMode,
	)

	m, err := migrate.New("file://./migrations", databaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize migration: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run cmd/migrate/main.go [up|down|steps N]")
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully.")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully.")

	case "steps":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run cmd/migrate/main.go steps <n>")
		}
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("Invalid step count: %v", err)
		}
		if err := m.Steps(n); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration steps failed: %v", err)
		}
		log.Printf("Migration steps (%d) applied successfully.", n)

	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}
