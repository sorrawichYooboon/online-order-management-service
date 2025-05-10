# Online Order Management Service

## Table of Contents

- [How to Set Up](#how-to-set-up)
- [How to Run the Server](#how-to-run-the-server)
- [Swagger](#swagger)
- [Project About](#project-about)

## How to Set Up

### Prerequisites

- Go (version 1.17 or higher)
- Docker

### 1. Run Docker Compose

To set up the PostgreSQL and Redis services, run the following command:

```base
docker-compose up -d
```

### 2. Check Docker Compose Run Correctly

```base
docker-compose ps
```

## How to Run the Server

### 1. Install Dependencies

```base
go mod tidy
```

### 2. Start the Server

```base
go run cmd/server/main.go
```

The server will start on port <b>8080</b>.

## Swagger (API Documentation)

You can access the Swagger documentation at:

```
http://localhost:8080/swagger/index.html
```

## Project About

### Project Structure

The project follows a modular structure adhering to Clean Architecture principles. Here’s a brief overview of the main directories:

![Clean Architecture DDD](https://storage.googleapis.com/bitloops-github-assets/Documentation%20Images/clean-architecture-and-ddd.png)

_Image Source: [Bitloops Documentation](https://bitloops.com/docs/bitloops-language/learning/software-architecture/clean-architecture)_

```
/online-order-management-service
├── cmd
│   ├── migrate
│   │   └── main.go                        # Entry point for database migrations
│   └── server
│       └── main.go                        # Main application entry point
├── config
│   └── config.go                          # Application configuration (env variables, struct loading)
├── docs
│   ├── docs.go                            # Swagger documentation metadata
│   ├── swagger.json                       # Generated Swagger JSON spec
│   └── swagger.yaml                       # Generated Swagger YAML spec
├── internal
│   ├── domain
│   │   ├── order.go                       # Domain entity: Order
│   │   └── order_item.go                 # Domain entity: OrderItem
│   ├── dto
│   │   └── order_dto.go                  # Data Transfer Objects (DTOs) for Order endpoints
│   ├── infrastructure
│   │   ├── database
│   │   │   ├── postgres
│   │   │   │   ├── order_item_repo.go    # PostgreSQL repository for order items
│   │   │   │   ├── order_repo.go         # PostgreSQL repository for orders
│   │   │   │   ├── pg_tx_manager.go      # PostgreSQL transaction manager implementation
│   │   │   │   └── postgres.go           # PostgreSQL connection setup
│   │   ├── delivery
│   │   │   └── rest
│   │   │       ├── handler
│   │   │       │   ├── constants.go      # Constants used in HTTP handlers
│   │   │       │   ├── interface.go      # Handler interfaces
│   │   │       │   ├── order_handler.go  # HTTP handler for orders
│   │   │       │   └── ping_handler.go   # Health check handler
│   │   │       └── router.go             # Echo router
├── internal
│   ├── repository
│   │   ├── order_item_repository.go      # Order item repository interface
│   │   ├── order_repository.go           # Order repository interface
│   │   └── pg_tx_manager.go              # Interface for managing DB transactions
│   ├── usecase
│   │   ├── order_usecase.go              # Business logic implementation for orders
│   │   ├── interface.go                  # Usecase interface
│   │   ├── constants.go                  # Constants for usecase logic
│   │   ├── request.go                    # Internal request models
│   │   └── response.go                   # Internal response models
├── logger
│   └── logger.go                         # Application-wide logging utility
├── migrations
│   ├── 1746689243_create_orders_and_order_items.up.sql     # Migration script to create tables
│   ├── 1746689243_create_orders_and_order_items.down.sql   # Migration rollback script
│   └── migration.go                      # Migration runner logic
├── pkg
│   ├── apperror
│   │   └── usecase_errors.go             # Centralized application error types for usecase
│   ├── httperror
│   │   └── http_errors.go                # HTTP error types for client responses
│   ├── response
│   │   ├── constants.go                  # Response code/message constants
│   │   └── response.go                   # JSON response wrapper helpers
│   ├── retry
│   │   └── retry.go                      # Retry mechanism
│   ├── validator
│   │   └── validator.go                  # Custom input validation logic
│   └── workers
│       └── worker_pool.go               # Worker pool implementation for concurrent tasks
├── .gitignore                           # Git ignore rules
├── docker-compose.yaml                  # Docker compose file for local development
├── go.mod                               # Go module definition
├── go.sum                               # Go dependency checksums
└── Makefile                             # Makefile for common dev/test/build commands
```

- **cmd/**: Contains entry points for the application.

  - **migrate/**: Contains the main entry point for running database migrations.
  - **server/**: Contains the main entry point for starting the application server.

- **config/**: Contains configuration files and environment variable loading logic. This includes the application configuration struct and loading functions to read environment variables.
- **docs/**: Contains the Swagger documentation files. The `docs.go` file contains metadata for the Swagger documentation, while `swagger.json` and `swagger.yaml` are generated specifications.
- **internal/**: Contains the core business logic of the application.
  - **domain/**: Contains domain entities and value objects.
  - **dto/**: Contains Data Transfer Objects (DTOs) used for communication to clients.
  - **infrastructure/**: Contains implementation details for external dependencies, such as database connections and HTTP handlers.
    - **database/**: Contains database-related code, including repository implementations and connection setup.
      - **postgres/**: Contains PostgreSQL-specific implementations for repositories and transaction management.
    - **delivery/**: Contains the delivery layer, which handles HTTP requests and responses.
      - **rest/**: Contains RESTful HTTP handlers and routing logic.
        - **handler/**: Contains HTTP handler implementations for various endpoints.
        - **router.go**: Sets up the Echo router
    - **middleware/**: Contains middleware functions for request processing.
  - **repository/**: Contains repository interfaces for data access.
  - **usecase/**: Contains use case implementations, which encapsulate business logic and interact with repositories.
    - **constants.go**: Contains constants used in the use case layer.
    - **request.go**: Contains request models for use cases.
    - **response.go**: Contains response models for use cases.
- **logger/**: Contains the logger implementation for the application. This includes the logger setup and configuration.
- **migrations/**: Contains database migration scripts and logic for running migrations.
- **pkg/**: Contains utility packages and helpers used throughout the application.
  - **apperror/**: Contains application error types and handling logic.
  - **httperror/**: Contains HTTP error types and handling logic.
  - **response/**: Contains response-related utilities, including JSON response wrappers.
  - **retry/**: Contains retry logic for transient errors.
  - **validator/**: Contains input validation logic.
  - **workers/**: Contains worker pool implementation for concurrent processing.
