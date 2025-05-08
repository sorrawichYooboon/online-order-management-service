package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sorrawichYooboon/online-order-management-service/config"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/database/postgres"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest"
	"github.com/sorrawichYooboon/online-order-management-service/migrations"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pg := postgres.Connect(cfg)
	migrations.RunMigrations(cfg)
	defer pg.Close(context.Background())

	e := rest.NewServer(cfg, pg)

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatalf("Shutting down server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(ctx)
}
