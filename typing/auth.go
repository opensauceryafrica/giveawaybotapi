package typing

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
