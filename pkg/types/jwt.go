package types

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	ID             string
	Name           string
	Email          string
	ProfilePicture string
	jwt.RegisteredClaims
}
