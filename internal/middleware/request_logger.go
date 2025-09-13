package middleware

import (
	"bookstore-api/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

// RequestLoggerMiddleware handles request logging
type RequestLoggerMiddleware struct{}

// NewRequestLoggerMiddleware creates a new request logger middleware
func NewRequestLoggerMiddleware() *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{}
}

// RequestLogger returns a request logging middleware
func (m *RequestLoggerMiddleware) RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Process request
		err := c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Log request details
		utils.LogRequest(c, duration, err)
		
		return err
	}
}

