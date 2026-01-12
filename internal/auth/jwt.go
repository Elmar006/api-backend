package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt/v5"
)

var TokenAuth *jwtauth.JWTAuth

func initJWWT() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production-12345"
	}

	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func CrateToken(userID uint, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * 30 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserId(ctx context.Context) (uint, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed to get claims from context: %v", err)
	}

	userIdFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id not found in JWT claims")
	}

	return uint(userIdFloat), nil
}

func Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(TokenAuth)
}

func Authenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator
}
