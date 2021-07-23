package auth

import (
	"fmt"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/golang-jwt/jwt"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
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

func CheckToken(tokenString, secret string) error {
	tok, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if s, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[DEBUG] %v\n", s)
			return nil, fmt.Errorf("Unexpected token signing: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("[ERROR] %v\n", err)
		return err
	}

	if !tok.Valid {
		log.Printf("[ERROR] %v\n", tok)
		return helpers.ErrUnauthorized
	}

	return nil
}
