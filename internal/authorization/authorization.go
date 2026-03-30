// Package authorization содержит код для авторизации пользователя
package authorization

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID *uuid.UUID
}

const TokenExp = time.Hour * 3

// BuildJWTString создание JWT
func BuildJWTString(userID *uuid.UUID) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET_KEY is not set")
	}
	if userID == nil {
		return "", fmt.Errorf("userID is nil")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserID получение userID по токену
func GetUserID(tokenString string) *uuid.UUID {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return nil
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Явно проверяем алгоритм, чтобы избежать некорректной подписи.
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil || token == nil || !token.Valid {
		return nil
	}
	return claims.UserID
}
