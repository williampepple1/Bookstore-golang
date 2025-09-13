package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware handles rate limiting
type RateLimitMiddleware struct{}

// NewRateLimitMiddleware creates a new rate limit middleware
func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{}
}

// RateLimit returns a rate limiting middleware
func (m *RateLimitMiddleware) RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100,                // Maximum number of requests
		Expiration: 1 * time.Minute,    // Time window
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Rate limit exceeded. Please try again later.",
			})
		},
	})
}

// StrictRateLimit returns a stricter rate limiting middleware for sensitive endpoints
func (m *RateLimitMiddleware) StrictRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        10,                 // Maximum number of requests
		Expiration: 1 * time.Minute,    // Time window
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   true,
				"message": "Rate limit exceeded for this endpoint. Please try again later.",
			})
		},
	})
}
