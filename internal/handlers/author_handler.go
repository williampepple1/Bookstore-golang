package handlers

import (
	"bookstore-api/internal/models"
	"bookstore-api/internal/services"
	"bookstore-api/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthorHandler handles author-related HTTP requests
type AuthorHandler struct {
	authorService *services.AuthorService
}

// NewAuthorHandler creates a new author handler
func NewAuthorHandler() *AuthorHandler {
	return &AuthorHandler{
		authorService: services.NewAuthorService(),
	}
}

// CreateAuthorRequest represents the request payload for creating an author
type CreateAuthorRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"required,email"`
	Biography string `json:"biography,omitempty"`
}

// UpdateAuthorRequest represents the request payload for updating an author
type UpdateAuthorRequest struct {
	Name      string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Email     string `json:"email,omitempty" validate:"omitempty,email"`
	Biography string `json:"biography,omitempty"`
}

// CreateAuthor creates a new author
func (h *AuthorHandler) CreateAuthor(c *fiber.Ctx) error {
	var req CreateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation failed",
			"details": err.Error(),
		})
	}

	author := &models.Author{
		Name:      req.Name,
		Email:     req.Email,
		Biography: req.Biography,
	}

	if err := h.authorService.CreateAuthor(author); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create author",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Author created successfully",
		"data":    author,
	})
}

// GetAuthor retrieves an author by ID
func (h *AuthorHandler) GetAuthor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid author ID",
			"details": err.Error(),
		})
	}

	author, err := h.authorService.GetAuthorByID(id)
	if err != nil {
		if err.Error() == "author not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get author",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Author retrieved successfully",
		"data":    author,
	})
}

// GetAllAuthors retrieves all authors with pagination
func (h *AuthorHandler) GetAllAuthors(c *fiber.Ctx) error {
	page, limit := getPaginationParams(c)

	authors, total, err := h.authorService.GetAllAuthors(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get authors",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Authors retrieved successfully",
		"data":    authors,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateAuthor updates an existing author
func (h *AuthorHandler) UpdateAuthor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid author ID",
			"details": err.Error(),
		})
	}

	var req UpdateAuthorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Validation failed",
			"details": err.Error(),
		})
	}

	updates := &models.Author{
		Name:      req.Name,
		Email:     req.Email,
		Biography: req.Biography,
	}

	if err := h.authorService.UpdateAuthor(id, updates); err != nil {
		if err.Error() == "author not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update author",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Author updated successfully",
	})
}

// DeleteAuthor deletes an author
func (h *AuthorHandler) DeleteAuthor(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid author ID",
			"details": err.Error(),
		})
	}

	if err := h.authorService.DeleteAuthor(id); err != nil {
		if err.Error() == "author not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Author not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete author",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Author deleted successfully",
	})
}

// SearchAuthors searches authors by name or email
func (h *AuthorHandler) SearchAuthors(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Search query is required",
		})
	}

	page, limit := getPaginationParams(c)

	authors, total, err := h.authorService.SearchAuthors(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to search authors",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Authors found successfully",
		"data":    authors,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// getPaginationParams extracts pagination parameters from the request
func getPaginationParams(c *fiber.Ctx) (int, int) {
	page := 1
	limit := 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}
