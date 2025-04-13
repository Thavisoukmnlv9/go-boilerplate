package auth

import (
	"errors"

	"github.com/Thavisoukmnlv9/go-boilerplate/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

// ValidateTokenAndExtractUserID validates the JWT token and extracts the userID
func ValidateTokenAndExtractRole(tokenString string) (string, string, error) {
	secret := config.LoadConfig().JWTSecret
	type Claims struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		println("token parse error:", err.Error())
		return "", "", errors.New("invalid token")
	}

	// Now claims is already in the correct type
	if claims.UserID == "" || claims.Role == "" {
		return "", "", errors.New("missing user_id or role in token")
	}

	return claims.UserID, claims.Role, nil
}
