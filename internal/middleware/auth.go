package middleware

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		rawToken := c.Get("Authorization")
		token := extractToken(rawToken)
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		userID, userRole, err := auth.ValidateTokenAndExtractRole(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		println("Extracted userRole:", userRole)
		c.Locals("userID", userID)
		c.Locals("userRole", userRole)
		return c.Next()
	}
}

func extractToken(header string) string {
	// Support "Bearer <token>"
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}
	return header
}
