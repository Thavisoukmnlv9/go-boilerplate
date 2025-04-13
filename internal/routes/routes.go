package routes

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/handlers"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/middleware"
	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all the routes for the app
func SetupRoutes(app *fiber.App, enforcer *casbin.Enforcer, authHandler *handlers.AuthHandler) {
	// Public Routes
	app.Post("/login", authHandler.Login)
	app.Post("/register", authHandler.Register)
	app.Post("/refresh", authHandler.RefreshToken)

	// Authenticated and Authorized Routes
	protected := app.Group("", middleware.AuthMiddleware(), middleware.CasbinMiddleware(enforcer))

	protected.Get("/admin", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome Admin"})
	})

	protected.Get("/user", func(c *fiber.Ctx) error {
		// Example user handler
		return c.JSON(fiber.Map{"message": "Welcome User"})
	})

	// Protected route example
	app.Get("/protected", middleware.AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Protected route accessed"})
	})
}
