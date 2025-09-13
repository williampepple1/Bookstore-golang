package handlers

import (
	"bookstore-api/internal/models"
	"bookstore-api/internal/services"
	"bookstore-api/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// BookHandler handles book-related HTTP requests
type BookHandler struct {
	bookService *services.BookService
}

// NewBookHandler creates a new book handler
func NewBookHandler() *BookHandler {
	return &BookHandler{
		bookService: services.NewBookService(),
	}
}

// CreateBookRequest represents the request payload for creating a book
type CreateBookRequest struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	ISBN        string     `json:"isbn" validate:"required,len=13"`
	Description string     `json:"description,omitempty"`
	Price       float64    `json:"price" validate:"required,min=0"`
	Stock       int        `json:"stock" validate:"min=0"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	AuthorID    string     `json:"author_id" validate:"required,uuid"`
	CategoryID  string     `json:"category_id" validate:"required,uuid"`
}

// UpdateBookRequest represents the request payload for updating a book
type UpdateBookRequest struct {
	Title       string     `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	ISBN        string     `json:"isbn,omitempty" validate:"omitempty,len=13"`
	Description string     `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int       `json:"stock,omitempty" validate:"omitempty,min=0"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	AuthorID    string     `json:"author_id,omitempty" validate:"omitempty,uuid"`
	CategoryID  string     `json:"category_id,omitempty" validate:"omitempty,uuid"`
}

// UpdateStockRequest represents the request payload for updating book stock
type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,min=0"`
}

// CreateBook creates a new book
func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	var req CreateBookRequest
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

	// Parse UUIDs
	authorID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid author ID",
			"details": err.Error(),
		})
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid category ID",
			"details": err.Error(),
		})
	}

	book := &models.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		PublishedAt: req.PublishedAt,
		AuthorID:    authorID,
		CategoryID:  categoryID,
	}

	if err := h.bookService.CreateBook(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create book",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Book created successfully",
		"data":    book,
	})
}

// GetBook retrieves a book by ID
func (h *BookHandler) GetBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid book ID",
			"details": err.Error(),
		})
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get book",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Book retrieved successfully",
		"data":    book,
	})
}

// GetAllBooks retrieves all books with pagination
func (h *BookHandler) GetAllBooks(c *fiber.Ctx) error {
	page, limit := getPaginationParams(c)

	books, total, err := h.bookService.GetAllBooks(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get books",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Books retrieved successfully",
		"data":    books,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateBook updates an existing book
func (h *BookHandler) UpdateBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid book ID",
			"details": err.Error(),
		})
	}

	var req UpdateBookRequest
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

	updates := &models.Book{
		Title:       req.Title,
		ISBN:        req.ISBN,
		Description: req.Description,
		PublishedAt: req.PublishedAt,
	}

	// Parse UUIDs if provided
	if req.AuthorID != "" {
		authorID, err := uuid.Parse(req.AuthorID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid author ID",
				"details": err.Error(),
			})
		}
		updates.AuthorID = authorID
	}

	if req.CategoryID != "" {
		categoryID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid category ID",
				"details": err.Error(),
			})
		}
		updates.CategoryID = categoryID
	}

	// Set price and stock if provided
	if req.Price != nil {
		updates.Price = *req.Price
	}
	if req.Stock != nil {
		updates.Stock = *req.Stock
	}

	if err := h.bookService.UpdateBook(id, updates); err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update book",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Book updated successfully",
	})
}

// DeleteBook deletes a book
func (h *BookHandler) DeleteBook(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid book ID",
			"details": err.Error(),
		})
	}

	if err := h.bookService.DeleteBook(id); err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete book",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Book deleted successfully",
	})
}

// GetBooksByAuthor retrieves books by author ID
func (h *BookHandler) GetBooksByAuthor(c *fiber.Ctx) error {
	authorIDStr := c.Params("authorId")
	authorID, err := uuid.Parse(authorIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid author ID",
			"details": err.Error(),
		})
	}

	page, limit := getPaginationParams(c)

	books, total, err := h.bookService.GetBooksByAuthor(authorID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get books by author",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Books retrieved successfully",
		"data":    books,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetBooksByCategory retrieves books by category ID
func (h *BookHandler) GetBooksByCategory(c *fiber.Ctx) error {
	categoryIDStr := c.Params("categoryId")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid category ID",
			"details": err.Error(),
		})
	}

	page, limit := getPaginationParams(c)

	books, total, err := h.bookService.GetBooksByCategory(categoryID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to get books by category",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Books retrieved successfully",
		"data":    books,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// SearchBooks searches books by title, ISBN, or description
func (h *BookHandler) SearchBooks(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Search query is required",
		})
	}

	page, limit := getPaginationParams(c)

	books, total, err := h.bookService.SearchBooks(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to search books",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Books found successfully",
		"data":    books,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UpdateBookStock updates book stock
func (h *BookHandler) UpdateBookStock(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid book ID",
			"details": err.Error(),
		})
	}

	var req UpdateStockRequest
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

	if err := h.bookService.UpdateBookStock(id, req.Stock); err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Book not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to update book stock",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Book stock updated successfully",
	})
}
