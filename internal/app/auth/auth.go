package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	Username string
	jwt.StandardClaims
}

type TokenWithExpTime struct {
	Token   string
	ExpTime time.Time
}

func CreateToken(username, secret string) (*TokenWithExpTime, error) {
	expTime := time.Now().Add(1 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	})

	stoken, err := token.SignedString([]byte(secret))
	return &TokenWithExpTime{Token: stoken, ExpTime: expTime}, err
}
