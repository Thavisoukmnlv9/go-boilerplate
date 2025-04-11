package services

import (
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/models"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/repositories"
)

// UserService defines the user service structure
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo}
}

// Example method
func (s *UserService) GetUser(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
