package server

import (
	"bookstore-api/internal/config"
	"bookstore-api/internal/handlers"
	"bookstore-api/internal/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// HTTPServer represents the HTTP server
type HTTPServer struct {
	app    *fiber.App
	config *config.Config
}

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(cfg *config.Config) *HTTPServer {
	// Create Fiber app with config
	app := fiber.New(fiber.Config{
		AppName: "Bookstore API v1.0.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 statuscode
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Initialize middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware()
	requestLoggerMiddleware := middleware.NewRequestLoggerMiddleware()

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
	}))
	app.Use(rateLimitMiddleware.RateLimit())
	app.Use(requestLoggerMiddleware.RequestLogger())

	return &HTTPServer{
		app:    app,
		config: cfg,
	}
}

// SetupRoutes configures all the routes
func (s *HTTPServer) SetupRoutes() {
	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware()
	rateLimitMiddleware := middleware.NewRateLimitMiddleware()

	// Health check routes
	healthHandler := handlers.NewHealthHandler()
	s.app.Get("/health", healthHandler.Health)
	s.app.Get("/ready", healthHandler.Ready)

	// API documentation
	docsHandler := handlers.NewDocsHandler()
	s.app.Get("/docs", docsHandler.GetAPIDocs)
	s.app.Get("/api/docs", docsHandler.GetAPIDocs)

	// API v1 routes
	api := s.app.Group("/api/v1")
	
	// Initialize handlers
	authorHandler := handlers.NewAuthorHandler()
	categoryHandler := handlers.NewCategoryHandler()
	bookHandler := handlers.NewBookHandler()
	
	// Author routes
	authors := api.Group("/authors")
	authors.Post("/", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), authorHandler.CreateAuthor)
	authors.Get("/", authorHandler.GetAllAuthors)
	authors.Get("/search", authorHandler.SearchAuthors)
	authors.Get("/:id", authorHandler.GetAuthor)
	authors.Put("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), authorHandler.UpdateAuthor)
	authors.Delete("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), authorHandler.DeleteAuthor)
	
	// Category routes
	categories := api.Group("/categories")
	categories.Post("/", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), categoryHandler.CreateCategory)
	categories.Get("/", categoryHandler.GetAllCategories)
	categories.Get("/search", categoryHandler.SearchCategories)
	categories.Get("/:id", categoryHandler.GetCategory)
	categories.Put("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), categoryHandler.UpdateCategory)
	categories.Delete("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), categoryHandler.DeleteCategory)
	
	// Book routes
	books := api.Group("/books")
	books.Post("/", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), bookHandler.CreateBook)
	books.Get("/", bookHandler.GetAllBooks)
	books.Get("/search", bookHandler.SearchBooks)
	books.Get("/author/:authorId", bookHandler.GetBooksByAuthor)
	books.Get("/category/:categoryId", bookHandler.GetBooksByCategory)
	books.Get("/:id", bookHandler.GetBook)
	books.Put("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), bookHandler.UpdateBook)
	books.Put("/:id/stock", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), bookHandler.UpdateBookStock)
	books.Delete("/:id", rateLimitMiddleware.StrictRateLimit(), authMiddleware.RequireAuth(), bookHandler.DeleteBook)

	// Root route
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Bookstore API",
			"version": "1.0.0",
			"status":  "running",
		})
	})
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	addr := s.config.Server.Host + ":" + s.config.Server.Port
	log.Printf("Starting HTTP server on %s", addr)
	return s.app.Listen(addr)
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown() error {
	log.Println("Shutting down HTTP server...")
	return s.app.Shutdown()
}

// GetApp returns the Fiber app instance (for testing)
func (s *HTTPServer) GetApp() *fiber.App {
	return s.app
}
