package handlers

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.authService.Register(input.Username, input.Password, input.Role); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Registration failed"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "User registered"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	accessToken, refreshToken, err := h.authService.Login(input.Username, input.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	return c.JSON(fiber.Map{"access_token": accessToken, "refresh_token": refreshToken})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	token, err := h.authService.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
	}
	return c.JSON(fiber.Map{"access_token": token})
}
