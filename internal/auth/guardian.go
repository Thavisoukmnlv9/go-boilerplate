package auth

import (
	"errors"

	"github.com/Thavisoukmnlv9/go-boilerplate/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateTokenAndExtractUserID validates the JWT token and extracts the userID
func ValidateTokenAndExtractRole(tokenString string) (string, string, error) {
	secret := config.LoadConfig().JWTSecret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	userID, ok1 := claims["user_id"].(string)
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		return "", "", errors.New("missing user_id or role in token")
	}

	return userID, role, nil
}
