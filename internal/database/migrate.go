package database

import (
	"bookstore-api/internal/config"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MigrationRecord represents a migration record in the database
type MigrationRecord struct {
	ID        int    `gorm:"primaryKey"`
	Version   string `gorm:"uniqueIndex;not null"`
	AppliedAt string `gorm:"not null"`
}

// TableName returns the table name for MigrationRecord
func (MigrationRecord) TableName() string {
	return "migrations"
}

// Migrate runs database migrations
func Migrate(cfg *config.Config) error {
	// Connect to database
	db, err := Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create migration tracking table
	if err := createMigrationTable(db); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Run manual SQL migrations first
	if err := runSQLMigrations(db, cfg); err != nil {
		return fmt.Errorf("failed to run SQL migrations: %w", err)
	}

	// Skip GORM auto-migrations since we're using manual SQL migrations
	// GORM auto-migrations can conflict with existing SQL schema
	log.Println("Skipping GORM auto-migrations (using manual SQL migrations)")

	log.Println("Database migrations completed successfully")
	return nil
}

// createMigrationTable creates the migration tracking table
func createMigrationTable(db *gorm.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			applied_at VARCHAR(255) NOT NULL
		)
	`).Error
}

// runSQLMigrations runs manual SQL migrations from the migrations directory
func runSQLMigrations(db *gorm.DB, cfg *config.Config) error {
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Println("No migrations directory found, skipping SQL migrations")
		return nil
	}

	// Get list of migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && !strings.HasPrefix(file.Name(), ".") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, file := range migrationFiles {
		version := strings.TrimSuffix(file, ".sql")

		// Skip if already applied
		if contains(appliedMigrations, version) {
			log.Printf("Migration %s already applied, skipping", version)
			continue
		}

		log.Printf("Applying migration: %s", version)

		// Read migration file
		filePath := filepath.Join(migrationsDir, file)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration
		if err := executeMigration(db, version, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		log.Printf("Successfully applied migration: %s", version)
	}

	return nil
}

// getAppliedMigrations returns a list of applied migration versions
func getAppliedMigrations(db *gorm.DB) ([]string, error) {
	var records []MigrationRecord
	if err := db.Find(&records).Error; err != nil {
		// If the table doesn't exist yet, return empty list
		if strings.Contains(err.Error(), "does not exist") {
			return []string{}, nil
		}
		return nil, err
	}

	var versions []string
	for _, record := range records {
		versions = append(versions, record.Version)
	}
	return versions, nil
}

// executeMigration executes a single migration and records it
func executeMigration(db *gorm.DB, version, content string) error {
	// Start transaction
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Execute migration SQL
	if err := tx.Exec(content).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Record migration
	record := MigrationRecord{
		Version:   version,
		AppliedAt: "now()",
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GetMigrationStatus returns the status of all migrations
func GetMigrationStatus(cfg *config.Config) ([]MigrationRecord, error) {
	db, err := Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	var records []MigrationRecord
	if err := db.Order("applied_at ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration status: %w", err)
	}

	return records, nil
}

// RollbackLastMigration rolls back the last applied migration
func RollbackLastMigration(cfg *config.Config) error {
	db, err := Connect(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the last applied migration
	var lastMigration MigrationRecord
	if err := db.Order("applied_at DESC").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	log.Printf("Rolling back migration: %s", lastMigration.Version)

	// For now, we'll just remove the migration record
	// In a production system, you'd want to implement proper rollback SQL
	if err := db.Delete(&lastMigration).Error; err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.Version, err)
	}

	log.Printf("Successfully rolled back migration: %s", lastMigration.Version)
	return nil
}

// ValidateMigrations checks if all migration files are properly formatted
func ValidateMigrations() error {
	migrationsDir := "migrations"

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Println("No migrations directory found")
		return nil
	}

	// Get list of migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Validate each migration file
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && !strings.HasPrefix(file.Name(), ".") {
			filePath := filepath.Join(migrationsDir, file.Name())
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
			}

			// Basic validation - check if file is not empty and contains SQL
			if len(strings.TrimSpace(string(content))) == 0 {
				return fmt.Errorf("migration file %s is empty", file.Name())
			}

			log.Printf("Migration file %s is valid", file.Name())
		}
	}

	log.Println("All migration files are valid")
	return nil
}

// Connect establishes a connection to the database
func Connect(cfg *config.Config) (*gorm.DB, error) {
	// First try to connect to the specific database
	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If database doesn't exist, try to create it
		if strings.Contains(err.Error(), "does not exist") {
			log.Println("Database does not exist, attempting to create it...")

			// Connect to postgres database to create the target database
			postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
				cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.SSLMode)

			postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{})
			if err != nil {
				return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
			}

			// Create the database
			if err := postgresDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName)).Error; err != nil {
				sqlDB, _ := postgresDB.DB()
				sqlDB.Close()
				return nil, fmt.Errorf("failed to create database %s: %w", cfg.Database.DBName, err)
			}

			sqlDB, _ := postgresDB.DB()
			sqlDB.Close()
			log.Printf("Database %s created successfully", cfg.Database.DBName)

			// Now try to connect to the newly created database
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, fmt.Errorf("failed to connect to newly created database: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connection established successfully")
	return db, nil
}
