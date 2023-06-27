package authorization

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/poggerr/go_shortener/internal/app/models"
	"github.com/poggerr/go_shortener/internal/app/storage"
	"github.com/poggerr/go_shortener/internal/encrypt"
	"github.com/poggerr/go_shortener/internal/logger"
	"os"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID *uuid.UUID
}

const TOKEN_EXP = time.Hour * 3

func BuildJWTString(uuid *uuid.UUID) (string, error) {

	var secretKey = os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: uuid,
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserID(tokenString string) string {
	var secretKey = os.Getenv("SECRET_KEY")
	claims := &Claims{}
	jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return claims.UserID.String()
}

func RegisterUser(strg *storage.Storage, user *models.User) {
	user.Pass = encrypt.Encrypt(user.Pass)
	id := uuid.New()
	err := strg.CreateUser(user.UserName, user.Pass, &id)
	if err != nil {
		logger.Initialize().Error(err)
	}
}
