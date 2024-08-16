package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
)

type JWTTokenGenerator struct {
	secret  string
	expires time.Duration
}

// Should implements TokenGenerator
var _ TokenGenerator = (*JWTTokenGenerator)(nil)

// This generator doesn't implements default logic of jwt, it follows DIP and SRP
// We create token generator as infrastructure and use it inside a usecase layer
// so we shouldn't depends on JWT.Claims or any other stuff
func NewJWTTokenGenerator(secret string, expires time.Duration) *JWTTokenGenerator {
	return &JWTTokenGenerator{
		secret:  secret,
		expires: expires,
	}
}

func (g *JWTTokenGenerator) Generate(payload interface{}) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(g.expires).Unix(),
	}

	if err := mapstructure.Decode(payload, claims); err != nil {
		return "", err
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

	claims, ok := payload.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to assert payload as jwt.MapClaims")
	}

	if err := mapstructure.Decode(claims, dest); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	return nil
}
