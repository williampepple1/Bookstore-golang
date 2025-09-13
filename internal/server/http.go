package server

import (
	"bookstore-api/internal/config"
	"bookstore-api/internal/handlers"
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

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	return &HTTPServer{
		app:    app,
		config: cfg,
	}
}

// SetupRoutes configures all the routes
func (s *HTTPServer) SetupRoutes() {
	// Health check routes
	healthHandler := handlers.NewHealthHandler()
	s.app.Get("/health", healthHandler.Health)
	s.app.Get("/ready", healthHandler.Ready)

	// API v1 routes
	api := s.app.Group("/api/v1")

	// Placeholder routes - will be implemented in next milestone
	api.Get("/books", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Books endpoint - coming soon",
			"status":  "pending",
		})
	})

	api.Get("/authors", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Authors endpoint - coming soon",
			"status":  "pending",
		})
	})

	api.Get("/categories", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Categories endpoint - coming soon",
			"status":  "pending",
		})
	})

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
