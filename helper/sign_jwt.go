package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/opensaucerers/giveawaybot/config"
	"github.com/opensaucerers/giveawaybot/typing"
)

// SignJWT signs a JWT with the given address
func SignJWT(id string, exp ...bool) (string, error) {
	expat := jwt.NewNumericDate(time.Now().UTC().Add(24 * time.Hour))
	if len(exp) > 0 && !exp[0] {
		expat = nil
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, typing.JWTClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			// expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: expat,
			Issuer:    config.Env.AppName,
		},
	}).SignedString([]byte(config.Env.JWTSecret))
}
