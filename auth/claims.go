package auth

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	ID string
	jwt.StandardClaims
}
