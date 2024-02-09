package util

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key []byte = []byte("ovalfi") // this should be loaded from environment when in production
	t   *jwt.Token
	s   string
)

// creates a signed jwt token
func CreateToken(subject string) (string, error) {
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "ovalfi",
		"sub": subject,
	})
	s, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return s, nil
}

// validates jwt token
func ValidateToken(s string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			os.Exit(1)
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}
	return claims, nil
}
