package handlers

import (
	"bookstore-api/internal/database"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns the health status of the application
func (h *HealthHandler) Health(c *fiber.Ctx) error {
	// Check if database is available
	if err := database.HealthCheck(); err != nil {
		// If database is not available, return partial health
		return c.JSON(fiber.Map{
			"status":  "degraded",
			"message": "Application running but database unavailable",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "healthy",
		"message": "All services are running",
	})
}

// Ready returns the readiness status of the application
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	// Check if database is ready
	if err := database.HealthCheck(); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"status":  "not ready",
			"message": "Database is not ready",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ready",
		"message": "Application is ready to serve requests",
	})
}
