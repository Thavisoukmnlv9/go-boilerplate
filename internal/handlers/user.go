package handlers

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/services"
	"github.com/gofiber/fiber/v2"
)

// UserHandler defines the user handler structure
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService}
}

// Example handler function
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
