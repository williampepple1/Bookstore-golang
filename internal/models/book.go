package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Book represents a book in the bookstore
type Book struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null;size:255" validate:"required,min=1,max=255"`
	ISBN        string         `json:"isbn" gorm:"uniqueIndex;not null;size:20" validate:"required,len=13"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"not null;type:decimal(10,2)" validate:"required,min=0"`
	Stock       int            `json:"stock" gorm:"not null;default:0" validate:"min=0"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Foreign Keys
	AuthorID   uuid.UUID `json:"author_id" gorm:"not null;type:uuid" validate:"required"`
	CategoryID uuid.UUID `json:"category_id" gorm:"not null;type:uuid" validate:"required"`

	// Relationships
	Author   Author   `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName returns the table name for the Book model
func (Book) TableName() string {
	return "books"
}

// BeforeCreate hook to generate UUID
func (b *Book) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
