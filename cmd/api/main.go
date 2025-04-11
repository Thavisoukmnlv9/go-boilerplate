package main

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/auth"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/config"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/handlers"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/middleware"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/repositories"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Failed to connect to database: ", err)
	}

	auth.Init()

	app := fiber.New()
	app.Use(middleware.RateLimitMiddleware())

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	app.Post("/register", authHandler.Register)
	app.Post("/login", authHandler.Login)
	app.Post("/refresh", authHandler.RefreshToken)

	app.Get("/protected", middleware.AuthMiddleware(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Protected route accessed"})
	})

	logrus.Info("Starting server on ", cfg.ServerPort)
	if err := app.Listen(cfg.ServerPort); err != nil {
		logrus.Fatal("Server failed to start: ", err)
	}
}
