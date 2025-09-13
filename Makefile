# Bookstore API Makefile

.PHONY: help build run test clean proto migrate migrate-status migrate-rollback migrate-validate migrate-up migrate-down dev-setup

# Default target
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests"
	@echo "  clean           - Clean build artifacts"
	@echo "  proto           - Generate protobuf files"
	@echo "  migrate         - Run database migrations"
	@echo "  migrate-status  - Check migration status"
	@echo "  migrate-rollback - Rollback last migration"
	@echo "  migrate-validate - Validate migration files"
	@echo "  migrate-up      - Alias for migrate"
	@echo "  migrate-down    - Alias for migrate-rollback"
	@echo "  dev-setup       - Setup development environment"

# Build the application
build:
	@echo "Building bookstore-api..."
	@go build -o bin/bookstore-api cmd/server/main.go

# Run the application
run:
	@echo "Running bookstore-api..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

# Database migrations
migrate:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go -action=migrate

migrate-status:
	@echo "Checking migration status..."
	@go run cmd/migrate/main.go -action=status

migrate-rollback:
	@echo "Rolling back last migration..."
	@go run cmd/migrate/main.go -action=rollback

migrate-validate:
	@echo "Validating migration files..."
	@go run cmd/migrate/main.go -action=validate

# Legacy migration commands (for compatibility)
migrate-up: migrate
migrate-down: migrate-rollback

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod tidy
	@go mod download
