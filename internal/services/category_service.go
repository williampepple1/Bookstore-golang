package services

import (
	"bookstore-api/internal/database"
	"bookstore-api/internal/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CategoryService handles category-related business logic
type CategoryService struct {
	db *gorm.DB
}

// NewCategoryService creates a new category service
func NewCategoryService() *CategoryService {
	return &CategoryService{
		db: database.GetDB(),
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(category *models.Category) error {
	if err := s.db.Create(category).Error; err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Books").First(&category, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// GetAllCategories retrieves all categories with pagination
func (s *CategoryService) GetAllCategories(page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	// Count total records
	if err := s.db.Model(&models.Category{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get categories with pagination
	if err := s.db.Preload("Books").Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	return categories, total, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id uuid.UUID, updates *models.Category) error {
	result := s.db.Model(&models.Category{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// DeleteCategory soft deletes a category
func (s *CategoryService) DeleteCategory(id uuid.UUID) error {
	result := s.db.Delete(&models.Category{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

// GetCategoryByName retrieves a category by name
func (s *CategoryService) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Books").First(&category, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

// SearchCategories searches categories by name or description
func (s *CategoryService) SearchCategories(query string, page, limit int) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	searchQuery := "%" + query + "%"

	// Count total records
	if err := s.db.Model(&models.Category{}).Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Search categories with pagination
	if err := s.db.Preload("Books").Where("name ILIKE ? OR description ILIKE ?", searchQuery, searchQuery).Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search categories: %w", err)
	}

	return categories, total, nil
}
