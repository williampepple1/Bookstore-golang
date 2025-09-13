package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Author represents an author in the bookstore
type Author struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string         `json:"name" gorm:"not null;size:255" validate:"required,min=2,max=255"`
	Email     string         `json:"email" gorm:"uniqueIndex:uni_authors_email;not null;size:255" validate:"required,email"`
	Biography string         `json:"biography" gorm:"type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Books []Book `json:"books,omitempty" gorm:"foreignKey:AuthorID"`
}

// TableName returns the table name for the Author model
func (Author) TableName() string {
	return "authors"
}

// BeforeCreate hook to generate UUID
func (a *Author) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
