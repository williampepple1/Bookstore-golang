package models

import (
	"gorm.io/gorm"
)

// AllModels returns a slice of all model structs for auto-migration
func AllModels() []interface{} {
	return []interface{}{
		&Author{},
		&Category{},
		&Book{},
	}
}

// Migrate runs database migrations for all models
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(AllModels()...)
}
