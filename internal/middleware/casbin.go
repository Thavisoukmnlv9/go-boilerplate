package middleware

import (
	"fmt"
	"strings"

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

		resource := strings.TrimRight(strings.Split(c.Path(), "?")[0], "/")
		action := c.Method()

		fmt.Printf("User Role: %s, Resource: %s, Action: %s\n", userRole, resource, action)

		// Check the policy enforcement
		ok, err := enforcer.Enforce(userRole, resource, action)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": fmt.Sprintf("Error enforcing policy: %v", err)})
		}
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}

		return c.Next()
	}
}
