package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"bookstore-api/internal/config"
	"bookstore-api/internal/database"
)

func main() {
	var (
		action = flag.String("action", "migrate", "Action to perform: migrate, status, rollback, validate")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	switch *action {
	case "migrate":
		if err := database.Migrate(cfg); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "status":
		migrations, err := database.GetMigrationStatus(cfg)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}

		if len(migrations) == 0 {
			fmt.Println("No migrations applied")
			return
		}

		fmt.Printf("Applied migrations (%d):\n", len(migrations))
		for _, migration := range migrations {
			fmt.Printf("  - %s (applied at: %s)\n", migration.Version, migration.AppliedAt)
		}

	case "rollback":
		if err := database.RollbackLastMigration(cfg); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Println("Rollback completed successfully")

	case "validate":
		if err := database.ValidateMigrations(); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		fmt.Println("All migration files are valid")

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: migrate, status, rollback, validate")
		os.Exit(1)
	}
}
