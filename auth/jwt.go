package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtHandler struct {
	secret []byte
	ExpirationTime time.Duration
}

func NewJwtHandler(secret []byte, duration time.Duration) JwtHandler {
	return JwtHandler{
		secret: secret,
		ExpirationTime: duration,
	}
}

func (jh *JwtHandler) GetNewToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jh.secret)
}

func (jh *JwtHandler) ParseClaims(tokenString string, emptyClaims jwt.Claims) (*jwt.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, emptyClaims, func(token *jwt.Token) (interface{}, error) {
		return jh.secret, nil
	})
	if err != nil {
		return nil, err
	}

	return &token.Claims, nil
}