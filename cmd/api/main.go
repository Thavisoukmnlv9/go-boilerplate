package main

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/auth"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/config"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/handlers"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/middleware"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/repositories"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/routes"
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

	// Connect to the database
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Failed to connect to database: ", err)
	}

	// Initialize app
	app := fiber.New()
	app.Use(middleware.RateLimitMiddleware())

	// Set up repositories, services, and handlers
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize Casbin enforcer
	enforcer, err := auth.InitializeCasbin(db)
	if err != nil {
		logrus.Fatal("Failed to initialize Casbin: ", err)
	}

	// Set up routes with Casbin authorization middleware
	routes.SetupRoutes(app, enforcer, authHandler)

	logrus.Info("Starting server on ", cfg.ServerPort)
	if err := app.Listen(cfg.ServerPort); err != nil {
		logrus.Fatal("Server failed to start: ", err)
	}
}
