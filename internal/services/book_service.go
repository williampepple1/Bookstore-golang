package services

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookService handles book-related business logic
type BookService struct {
	db *gorm.DB
}

// NewBookService creates a new book service
func NewBookService() *BookService {
	return &BookService{
		db: database.GetDB(),
	}
}

// CreateBook creates a new book
func (s *BookService) CreateBook(book *models.Book) error {
	// Validate that author and category exist
	if err := s.validateAuthorAndCategory(book.AuthorID, book.CategoryID); err != nil {
		return err
	}

	if err := s.db.Create(book).Error; err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}
	return nil
}

// GetBookByID retrieves a book by ID
func (s *BookService) GetBookByID(id uuid.UUID) (*models.Book, error) {
	var book models.Book
	if err := s.db.Preload("Author").Preload("Category").First(&book, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("failed to get book: %w", err)
	}
	return &book, nil
}

// GetAllBooks retrieves all books with pagination
func (s *BookService) GetAllBooks(page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Count total records
	if err := s.db.Model(&models.Book{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get books with pagination
	if err := s.db.Preload("Author").Preload("Category").Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get books: %w", err)
	}

	return books, total, nil
}

// UpdateBook updates an existing book
func (s *BookService) UpdateBook(id uuid.UUID, updates *models.Book) error {
	// If updating author or category, validate they exist
	if updates.AuthorID != uuid.Nil || updates.CategoryID != uuid.Nil {
		// Get current book to check existing values
		var currentBook models.Book
		if err := s.db.First(&currentBook, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("book not found")
			}
			return fmt.Errorf("failed to get book: %w", err)
		}

		authorID := currentBook.AuthorID
		categoryID := currentBook.CategoryID

		if updates.AuthorID != uuid.Nil {
			authorID = updates.AuthorID
		}
		if updates.CategoryID != uuid.Nil {
			categoryID = updates.CategoryID
		}

		if err := s.validateAuthorAndCategory(authorID, categoryID); err != nil {
			return err
		}
	}

	result := s.db.Model(&models.Book{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update book: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("book not found")
	}
	return nil
}

// DeleteBook soft deletes a book
func (s *BookService) DeleteBook(id uuid.UUID) error {
	result := s.db.Delete(&models.Book{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete book: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("book not found")
	}
	return nil
}

// GetBooksByAuthor retrieves books by author ID
func (s *BookService) GetBooksByAuthor(authorID uuid.UUID, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Count total records
	if err := s.db.Model(&models.Book{}).Where("author_id = ?", authorID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get books with pagination
	if err := s.db.Preload("Author").Preload("Category").Where("author_id = ?", authorID).Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get books: %w", err)
	}

	return books, total, nil
}

// GetBooksByCategory retrieves books by category ID
func (s *BookService) GetBooksByCategory(categoryID uuid.UUID, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Count total records
	if err := s.db.Model(&models.Book{}).Where("category_id = ?", categoryID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get books with pagination
	if err := s.db.Preload("Author").Preload("Category").Where("category_id = ?", categoryID).Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get books: %w", err)
	}

	return books, total, nil
}

// SearchBooks searches books by title, ISBN, or description
func (s *BookService) SearchBooks(query string, page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Book{}).Where("title ILIKE ? OR isbn ILIKE ? OR description ILIKE ?", searchQuery, searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count books: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Search books with pagination
	if err := s.db.Preload("Author").Preload("Category").Where("title ILIKE ? OR isbn ILIKE ? OR description ILIKE ?", searchQuery, searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search books: %w", err)
	}

	return books, total, nil
}

// UpdateBookStock updates book stock
func (s *BookService) UpdateBookStock(id uuid.UUID, newStock int) error {
	if newStock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}

	result := s.db.Model(&models.Book{}).Where("id = ?", id).Update("stock", newStock)
	if result.Error != nil {
		return fmt.Errorf("failed to update book stock: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("book not found")
	}
	return nil
}

// validateAuthorAndCategory validates that author and category exist
func (s *BookService) validateAuthorAndCategory(authorID, categoryID uuid.UUID) error {
	// Check if author exists
	var authorCount int64
	if err := s.db.Model(&models.Author{}).Where("id = ?", authorID).Count(&authorCount).Error; err != nil {
		return fmt.Errorf("failed to validate author: %w", err)
	}
	if authorCount == 0 {
		return fmt.Errorf("author not found")
	}

	// Check if category exists
	var categoryCount int64
	if err := s.db.Model(&models.Category{}).Where("id = ?", categoryID).Count(&categoryCount).Error; err != nil {
		return fmt.Errorf("failed to validate category: %w", err)
	}
	if categoryCount == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
