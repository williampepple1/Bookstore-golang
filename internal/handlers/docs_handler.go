package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// DocsHandler handles API documentation
type DocsHandler struct{}

// NewDocsHandler creates a new docs handler
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// GetAPIDocs returns API documentation
func (h *DocsHandler) GetAPIDocs(c *fiber.Ctx) error {
	docs := fiber.Map{
		"title":       "Bookstore API",
		"version":     "1.0.0",
		"description": "A comprehensive bookstore management API",
		"base_url":    "http://localhost:8080/api/v1",
		"endpoints": fiber.Map{
			"authors": fiber.Map{
				"description": "Author management endpoints",
				"endpoints": []fiber.Map{
					{
						"method":      "GET",
						"path":        "/authors",
						"description": "List all authors with pagination",
						"parameters":  []string{"page", "limit"},
						"response":    "List of authors with pagination info",
					},
					{
						"method":      "POST",
						"path":        "/authors",
						"description": "Create a new author",
						"body":        "Author data (name, email, biography)",
						"response":    "Created author object",
					},
					{
						"method":      "GET",
						"path":        "/authors/:id",
						"description": "Get author by ID",
						"parameters":  []string{"id (UUID)"},
						"response":    "Author object with books",
					},
					{
						"method":      "PUT",
						"path":        "/authors/:id",
						"description": "Update author",
						"parameters":  []string{"id (UUID)"},
						"body":        "Updated author data",
						"response":    "Success message",
					},
					{
						"method":      "DELETE",
						"path":        "/authors/:id",
						"description": "Delete author",
						"parameters":  []string{"id (UUID)"},
						"response":    "Success message",
					},
					{
						"method":      "GET",
						"path":        "/authors/search",
						"description": "Search authors",
						"parameters":  []string{"q (query string)"},
						"response":    "List of matching authors",
					},
				},
			},
			"categories": fiber.Map{
				"description": "Category management endpoints",
				"endpoints": []fiber.Map{
					{
						"method":      "GET",
						"path":        "/categories",
						"description": "List all categories with pagination",
						"parameters":  []string{"page", "limit"},
						"response":    "List of categories with pagination info",
					},
					{
						"method":      "POST",
						"path":        "/categories",
						"description": "Create a new category",
						"body":        "Category data (name, description)",
						"response":    "Created category object",
					},
					{
						"method":      "GET",
						"path":        "/categories/:id",
						"description": "Get category by ID",
						"parameters":  []string{"id (UUID)"},
						"response":    "Category object with books",
					},
					{
						"method":      "PUT",
						"path":        "/categories/:id",
						"description": "Update category",
						"parameters":  []string{"id (UUID)"},
						"body":        "Updated category data",
						"response":    "Success message",
					},
					{
						"method":      "DELETE",
						"path":        "/categories/:id",
						"description": "Delete category",
						"parameters":  []string{"id (UUID)"},
						"response":    "Success message",
					},
					{
						"method":      "GET",
						"path":        "/categories/search",
						"description": "Search categories",
						"parameters":  []string{"q (query string)"},
						"response":    "List of matching categories",
					},
				},
			},
			"books": fiber.Map{
				"description": "Book management endpoints",
				"endpoints": []fiber.Map{
					{
						"method":      "GET",
						"path":        "/books",
						"description": "List all books with pagination",
						"parameters":  []string{"page", "limit"},
						"response":    "List of books with pagination info",
					},
					{
						"method":      "POST",
						"path":        "/books",
						"description": "Create a new book",
						"body":        "Book data (title, isbn, description, price, stock, author_id, category_id)",
						"response":    "Created book object",
					},
					{
						"method":      "GET",
						"path":        "/books/:id",
						"description": "Get book by ID",
						"parameters":  []string{"id (UUID)"},
						"response":    "Book object with author and category",
					},
					{
						"method":      "PUT",
						"path":        "/books/:id",
						"description": "Update book",
						"parameters":  []string{"id (UUID)"},
						"body":        "Updated book data",
						"response":    "Success message",
					},
					{
						"method":      "DELETE",
						"path":        "/books/:id",
						"description": "Delete book",
						"parameters":  []string{"id (UUID)"},
						"response":    "Success message",
					},
					{
						"method":      "GET",
						"path":        "/books/search",
						"description": "Search books",
						"parameters":  []string{"q (query string)"},
						"response":    "List of matching books",
					},
					{
						"method":      "GET",
						"path":        "/books/author/:authorId",
						"description": "Get books by author",
						"parameters":  []string{"authorId (UUID)", "page", "limit"},
						"response":    "List of books by author",
					},
					{
						"method":      "GET",
						"path":        "/books/category/:categoryId",
						"description": "Get books by category",
						"parameters":  []string{"categoryId (UUID)", "page", "limit"},
						"response":    "List of books by category",
					},
					{
						"method":      "PUT",
						"path":        "/books/:id/stock",
						"description": "Update book stock",
						"parameters":  []string{"id (UUID)"},
						"body":        "Stock data (stock: number)",
						"response":    "Success message",
					},
				},
			},
			"health": fiber.Map{
				"description": "Health check endpoints",
				"endpoints": []fiber.Map{
					{
						"method":      "GET",
						"path":        "/health",
						"description": "Check application health",
						"response":    "Health status",
					},
					{
						"method":      "GET",
						"path":        "/ready",
						"description": "Check application readiness",
						"response":    "Readiness status",
					},
				},
			},
		},
		"authentication": fiber.Map{
			"type":        "Bearer Token",
			"description": "Include 'Authorization: Bearer <token>' header for protected endpoints",
			"note":        "Currently using placeholder authentication",
		},
		"pagination": fiber.Map{
			"parameters": []string{"page (default: 1)", "limit (default: 10, max: 100)"},
			"response":   "Includes pagination info with total, total_pages, page, limit",
		},
		"error_format": fiber.Map{
			"structure": fiber.Map{
				"error":   "boolean",
				"message": "string",
				"details": "string (optional)",
			},
			"example": fiber.Map{
				"error":   true,
				"message": "Validation failed",
				"details": "Name is required; Email must be a valid email address",
			},
		},
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "API documentation retrieved successfully",
		"data":    docs,
	})
}
