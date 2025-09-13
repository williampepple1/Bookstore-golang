package services

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthorService handles author-related business logic
type AuthorService struct {
	db *gorm.DB
}

// NewAuthorService creates a new author service
func NewAuthorService() *AuthorService {
	return &AuthorService{
		db: database.GetDB(),
	}
}

// CreateAuthor creates a new author
func (s *AuthorService) CreateAuthor(author *models.Author) error {
	if err := s.db.Create(author).Error; err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	return nil
}

// GetAuthorByID retrieves an author by ID
func (s *AuthorService) GetAuthorByID(id uuid.UUID) (*models.Author, error) {
	var author models.Author
	if err := s.db.Preload("Books").First(&author, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	return &author, nil
}

// GetAllAuthors retrieves all authors with pagination
func (s *AuthorService) GetAllAuthors(page, limit int) ([]models.Author, int64, error) {
	var authors []models.Author
	var total int64

	// Count total records
	if err := s.db.Model(&models.Author{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count authors: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get authors with pagination
	if err := s.db.Preload("Books").Offset(offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get authors: %w", err)
	}

	return authors, total, nil
}

// UpdateAuthor updates an existing author
func (s *AuthorService) UpdateAuthor(id uuid.UUID, updates *models.Author) error {
	result := s.db.Model(&models.Author{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update author: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("author not found")
	}
	return nil
}

// DeleteAuthor soft deletes an author
func (s *AuthorService) DeleteAuthor(id uuid.UUID) error {
	result := s.db.Delete(&models.Author{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete author: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("author not found")
	}
	return nil
}

// GetAuthorByEmail retrieves an author by email
func (s *AuthorService) GetAuthorByEmail(email string) (*models.Author, error) {
	var author models.Author
	if err := s.db.Preload("Books").First(&author, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("author not found")
		}
		return nil, fmt.Errorf("failed to get author: %w", err)
	}
	return &author, nil
}

// SearchAuthors searches authors by name or email
func (s *AuthorService) SearchAuthors(query string, page, limit int) ([]models.Author, int64, error) {
	var authors []models.Author
	var total int64

	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Author{}).Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count authors: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Search authors with pagination
	if err := s.db.Preload("Books").Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&authors).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search authors: %w", err)
	}

	return authors, total, nil
}
