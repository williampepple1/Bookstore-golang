package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"bookstore-api/internal/config"
	"bookstore-api/internal/database"
	"bookstore-api/internal/grpc"
	"bookstore-api/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Bookstore API server on port %s", cfg.Server.Port)
	log.Printf("Database: %s", cfg.Database.Host)

	// Initialize database connection using singleton pattern
	if err := database.InitializeDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Validate migration files before running
	if err := database.ValidateMigrations(); err != nil {
		log.Fatalf("Migration validation failed: %v", err)
	}

	// Run database migrations
	if err := database.Migrate(cfg); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Log migration status
	migrations, err := database.GetMigrationStatus(cfg)
	if err != nil {
		log.Printf("Warning: Failed to get migration status: %v", err)
	} else {
		log.Printf("Applied migrations: %d", len(migrations))
		for _, migration := range migrations {
			log.Printf("  - %s (applied at: %s)", migration.Version, migration.AppliedAt)
		}
	}

	log.Printf("Database connection established successfully")

	// Initialize servers
	httpServer := server.NewHTTPServer(cfg)
	httpServer.SetupRoutes()

	grpcServer := grpc.NewGRPCServer()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		if err := httpServer.Shutdown(); err != nil {
			log.Printf("Error shutting down HTTP server: %v", err)
		}
		// gRPC server will be stopped when the process exits
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		os.Exit(0)
	}()

	log.Println("Starting servers...")

	// Start HTTP server in goroutine
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start gRPC server in goroutine
	go func() {
		if err := grpcServer.Start(cfg); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Keep the main goroutine alive
	select {}
}
