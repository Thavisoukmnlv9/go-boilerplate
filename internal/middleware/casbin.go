package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

// CasbinMiddleware checks if the user has the right permissionsfunc CasbinMiddleware(enforcer *casbin.Enforcer) fiber.Handler {
func CasbinMiddleware(enforcer *casbin.Enforcer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Role not found in context"})
		}

		resource := c.Path()
		action := c.Method()

		// Enforce using the role
		ok, err := enforcer.Enforce(userRole, resource, action)
		if err != nil || !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}

		return c.Next()
	}
}
