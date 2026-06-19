// Package utils provides shared utilities such as JWT token handling
package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the custom claims embedded in access tokens
type JWTClaims struct {
	UserID string   `json:"user_id"`
	Login  string   `json:"login"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a signed JWT token for the given user identity
func GenerateJWT(userID, login string, roles []string, secret string, expires time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Login:  login,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseJWT validates and extracts claims from a JWT token string
func ParseJWT(tokenStr, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
