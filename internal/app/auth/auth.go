package auth

import "github.com/golang-jwt/jwt"

type CustomClaims struct {
	Username string
	jwt.StandardClaims
}
