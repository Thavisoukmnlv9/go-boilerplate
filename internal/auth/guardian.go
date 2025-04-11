package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Thavisoukmnlv9/go-boilerplate/internal/utils"
	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/auth/strategies/token"
	"github.com/shaj13/go-guardian/store"
)

var Strategy auth.Strategy
var Cache store.Cache

func Init() {
	// Cache with a 10 minute expiration window
	Cache = store.NewFIFO(context.Background(), 10*time.Minute)
	Strategy = token.New(validateToken, Cache)
}

// Updated function signature to include *http.Request
func validateToken(ctx context.Context, req *http.Request, rawToken string) (auth.Info, error) {
	claims, err := utils.ParseJWT(rawToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Check for token expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	// Return the authenticated user with claims.
	return auth.NewDefaultUser(claims.UserID, claims.UserID, nil, map[string][]string{
		"role": {claims.Role},
	}), nil
}
