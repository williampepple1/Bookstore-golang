# Bookstore API Makefile

.PHONY: help build run test clean proto migrate-up migrate-down

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  proto      - Generate protobuf files"
	@echo "  migrate-up - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"

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

# Database migrations (placeholder - will be implemented later)
migrate-up:
	@echo "Running database migrations up..."
	@echo "TODO: Implement database migrations"

migrate-down:
	@echo "Running database migrations down..."
	@echo "TODO: Implement database migrations"

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@go mod tidy
	@go mod download
