package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sorrawichYooboon/online-order-management-service/config"
	_ "github.com/sorrawichYooboon/online-order-management-service/docs"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/database/postgres"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest"
	"github.com/sorrawichYooboon/online-order-management-service/internal/infrastructure/delivery/rest/handler"
	"github.com/sorrawichYooboon/online-order-management-service/internal/usecase"
	"github.com/sorrawichYooboon/online-order-management-service/migrations"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/validator"
	"github.com/sorrawichYooboon/online-order-management-service/pkg/workers"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Online Order Management Service API
// @version 1.0
// @description KKP Interview - API for managing online orders and order items
// @host localhost:8080
// @BasePath /

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	pg := postgres.Connect(cfg)
	defer pg.Close()

	migrations.RunMigrations(cfg)

	pool := workers.NewWorkerPool(20, 1000)
	pool.Start()
	defer pool.Stop()

	e := echo.New()
	e.Validator = validator.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("5M"))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	pingHandler := handler.NewHealthHandler()

	pgTxManager := postgres.NewTxManager(pg)
	orderItemRepo := postgres.NewOrderItemRepository(pg)
	orderRepo := postgres.NewOrderRepository(pg)

	orderUsecase := usecase.NewOrderUsecase(pgTxManager, orderRepo, orderItemRepo, pool)

	orderHandler := handler.NewOrderHandler(orderUsecase)

	rest.NewServer(e, cfg, pingHandler, orderHandler)

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
