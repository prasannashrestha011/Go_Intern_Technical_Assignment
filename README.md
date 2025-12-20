# Go Intern Backend Project

A mini User & Order Management API built in Go demonstrating multiple HTTP routers, JWT authentication, PostgreSQL integration, simple concurrency, and clean architecture.

## Objective

The goal of this project is to showcase backend development skills in Go, including:

- Building REST APIs with different routers (Gin, Chi, Gorilla Mux)
- Implementing authentication and JWT
- Database integration with PostgreSQL using GORM
- Structured logging with Zap
- Configuration management with Viper
- Simple implementation of go concurrency and channels to understand its use cases.

## Project Structure

```
go-intern-task/
│
├── cmd/
│   ├── auth/          # Auth service (Gin)
│   ├── user/          # User service (Chi)
│   └── order/         # Order service (Gorilla Mux)
│
├── internal/
│   ├── config/        # Viper config
│   ├── database/      # GORM models and DB connection
│   ├──handlers/       # app routers
|   ├── logger/        # Zap logger
│   ├──
|   ├── middlewares/    # JWT and HTTP middlewares
|                    ├──gin_middleware
|                    ├──chi_middleware (compatible with mux router)
│   ├── models/        # DB models
│   ├── repository/    # DB operations
|   ├──schema          # app DTOs
│   ├── services/      # Business logic
|   ├──
│   └── utils/         # app helpers functions
|
│
├── migrations/        # Database migration scripts
├── go.sum
├── go.mod
└── README.md
```

## Services

### Auth Service (Gin)

- Handles user login and JWT generation
- Routes:
  - POST /auth/login
  - POST /auth/refresh
  - GET /auth/validate (protected)
  - GET /auth/profile (protected)

### User Service (Chi)

- CRUD operations for users
- JWT-protected endpoints
- Routes:
  - POST /create
  - GET /users/ <-- Get all users
  - GET /users/{id}
  - PUT /users/{id}
  - DELETE /users/{id}

### Order Service (Gorilla Mux)

- Order creation and retrieval
- JWT-protected
- Handles user-order relationship
- Routes:

  - Get /orders
  - POST /orders?page=nth&pageSize=n
  - GET /orders/{id}
  - PUT /orders/{id}
  - GET /users/{id}/orders?page=nth&pageSize=n

  ## Configuration

- Configuration is managed with Viper.
- Uses a `.yml` file for environment-specific settings.
- config.yml file should be placed in internal\config\ folder

```
server:
  port:
    auth: port1
    user: port2
    order: port3
  env: DEV || PROD

database:
  url: postgres_db_url

jwt:
  secret: jwt_secret_key
  accessExpiry: exp_in_hours
  refreshExpiry: exp_in_hours

```

## Migration configuration

This Project uses goose library for migrating schema into the database
I have created a makefile for shortening the goose commands

- Creating migration file

```
make migrate-create name=migration_name
```

- Applying all pending migration to the database

```
make migrate-up
```

- Rollback Last applied migration

```
make migrate-down
```

- Current migration status

```
make migrate-status
```

## Tradeoffs

- During containerization, I faced an issue with using YAML-based configuration (Viper) alongside Docker Compose.

  The application reads all configuration values (including database credentials) from a .yml file using Viper. However, Docker Compose is designed to inject configuration more naturally through environment variables (via .env files). Unlike .env, YAML files cannot be directly imported or interpolated inside docker-compose.yml.

  Because of this mismatch:

  Database credentials had to be hardcoded inside the YAML config

  Docker Compose could not dynamically override these values

  This prevented a clean, production-ready Docker setup at the final stage

  As a result, the project could not be fully dockerized in an ideal way within the given time.
