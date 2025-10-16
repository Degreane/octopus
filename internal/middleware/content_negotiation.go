// Package middleware provides HTTP middleware functions for the application.
//
// This package contains middleware components that can be used to process
// HTTP requests and responses in a Fiber web application.
package middleware

import "github.com/gofiber/fiber/v2"

// ContentNegotiation returns a Fiber middleware handler that determines the
// appropriate response format based on request headers.
//
// The middleware checks for:
// 1. HTMX requests via the HX-Request header
// 2. JSON requests via the Accept: application/json header
//
// It sets a "render" value in the Fiber context locals that can be either:
// - "html" for HTML responses (default)
// - "json" for JSON responses
//
// Usage:
//
//	app := fiber.New()
//	app.Use(middleware.ContentNegotiation())
//
// The downstream handlers can check the render format using:
//
//	format := c.Locals("render").(string)
func ContentNegotiation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for HTMX request
		if c.Get("HX-Request") != "" {
			c.Locals("render", "html")
			return c.Next()
		}

		// Check Accept header
		accept := c.Get("Accept")
		if accept == "application/json" {
			c.Locals("render", "json")
		} else {
			c.Locals("render", "html")
		}

		return c.Next()
	}
}
