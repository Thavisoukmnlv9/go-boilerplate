package services

import (
	"errors"
	"fmt"

	"github.com/Thavisoukmnlv9/go-boilerplate/internal/models"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/repositories"
	"github.com/Thavisoukmnlv9/go-boilerplate/internal/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo}
}

func (s *AuthService) Register(username, password, role string) error {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	user := &models.User{
		Username: username,
		Password: hash,
		Role:     role,
	}
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", "", err
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}
	// Convert user.ID (uint) to string before passing it to GenerateAccessToken
	accessToken, err := utils.GenerateAccessToken(fmt.Sprintf("%d", user.ID), user.Role)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := utils.ParseJWT(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Check if claims.Valid() returns an error or nil
	if err := claims.Valid(); err != nil {
		return "", errors.New("invalid refresh token")
	}

	return utils.GenerateAccessToken(claims.UserID, claims.Role)
}
