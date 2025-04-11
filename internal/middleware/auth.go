package middleware

import (
	"net/http"

	"github.com/Thavisoukmnlv9/go-boilerplate/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve the Authorization header from Fiber's context.
		raw := c.Get("Authorization")
		if raw == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		// Create a minimal HTTP request and set the Authorization header.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
		req.Header.Set("Authorization", raw)

		// Call Authenticate with only the context and the constructed HTTP request.
		_, err = auth.Strategy.Authenticate(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		return c.Next()
	}
}
