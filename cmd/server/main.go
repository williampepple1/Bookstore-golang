package main

import (
	"log"
	"os"

	"bookstore-api/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Bookstore API server on port %s", cfg.Server.Port)
	log.Printf("Database: %s", cfg.Database.Host)

	// TODO: Initialize database connection
	// TODO: Initialize HTTP server (Fiber)
	// TODO: Initialize gRPC server
	// TODO: Start servers

	log.Println("Server started successfully")
}
