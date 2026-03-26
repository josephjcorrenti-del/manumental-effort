package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	signingKey   string
	expiryMinute int
}

func NewTokenManager(signingKey string, expiryMinute int) *TokenManager {
	return &TokenManager{
		signingKey:   signingKey,
		expiryMinute: expiryMinute,
	}
}

func (m *TokenManager) CreateToken(userID string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Duration(m.expiryMinute) * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     expiresAt.Unix(),
	})

	signed, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (m *TokenManager) ParseToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userIDValue, ok := claims["user_id"]
	if !ok {
		return "", fmt.Errorf("missing user_id claim")
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid user_id claim")
	}

	return userID, nil
}
