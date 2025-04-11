package lib

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/immxrtalbeast/flagman-backend/internal/domain"
)

func NewToken(user *domain.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Email
	claims["fullname"] = user.FullName
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}
