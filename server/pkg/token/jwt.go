package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	secret  string
	expires time.Duration
}

func NewJWTTokenGenerator(secret string, expires time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secret:  secret,
		expires: expires,
	}
}

func (g *JWTTokenGenerator) Generate(payload interface{}) (string, error) {
	claims := jwt.MapClaims{
		"payload": payload,
		"exp":     time.Now().Add(g.expires).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (g *JWTTokenGenerator) Verify(token string) (interface{}, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(g.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return t.Claims, nil
}

func (g *JWTTokenGenerator) VerifyAndScan(token string, dest interface{}) error {
	payload, err := g.Verify(token)
	if err != nil {
		return err
	}
	dest = payload
	return nil
}
