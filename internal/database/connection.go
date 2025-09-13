package database

import (
	"bookstore-api/internal/config"
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB returns the singleton database connection
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitializeDB first.")
	}
	return db
}

// InitializeDB initializes the database connection
func InitializeDB(cfg *config.Config) error {
	var err error
	once.Do(func() {
		db, err = Connect(cfg)
		if err != nil {
			err = fmt.Errorf("failed to initialize database: %w", err)
			return
		}
	})
	return err
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck checks if the database connection is healthy
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Ping()
}
