package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/the-NZA/acg-nikolaev/internal/app/helpers"
)

const (
	tokenTTL  = 2 // Number of living hours for token
	renewTime = 5 // Number of minutes needed to start token update process
)

// CreateToken generate new token with passed params
func CreateToken(username, secret string) (string, time.Time, error) {
	expTime := time.Now().Add(tokenTTL * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expTime.Unix(),
	})

	stoken, err := token.SignedString([]byte(secret))
	return stoken, expTime, err
}

// CheckToken verify tokenString with given secret and return bool and error
// bool – signals that token may be updated
func CheckToken(tokenString, secret string) (bool, error) {
	tok, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected token signing: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	if !tok.Valid {
		return false, helpers.ErrUnauthorized
	}

	// Check if exp time less than
	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		exp := int64(claims["exp"].(float64))

		expTm := time.Unix(exp, 0)
		curTm := time.Now()
		dur := expTm.Sub(curTm)

		if dur.Minutes() < renewTime {
			return true, nil
		}
	}

	return false, nil
}

// UpdateToken generates and  returns new token, expTime and error
func UpdateToken(oldToken, secret string) (string, time.Time, error) {
	oldTokParsed, err := jwt.Parse(oldToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected token signing: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return "", time.Time{}, err
	}

	oldClaimg, ok := oldTokParsed.Claims.(jwt.MapClaims)
	if !ok {
		return "", time.Time{}, fmt.Errorf("Error during updating token")
	}

	username, ok := oldClaimg["username"].(string)
	if !ok {
		return "", time.Time{}, fmt.Errorf("Can't extract one or more claims fields")
	}

	return CreateToken(username, secret)

}
