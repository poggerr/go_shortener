// Package authorization содержит код для авторизации пользователя
package authorization

import (
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
func BuildJWTString(uuid *uuid.UUID) (string, error) {

	var secretKey = os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: uuid,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetUserID получение userID по токену
func GetUserID(tokenString string) string {
	//var secretKey = os.Getenv("SECRET_KEY")
	var secretKey = "scdcsdc,HVJHVCAJscdJccdsJVDVJDvqwe[p[;cqsc09cah989h"
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return claims.UserID.String()
}
