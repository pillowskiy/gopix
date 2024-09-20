package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
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

	if err := g.decode(payload, claims); err != nil {
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
		return nil, errors.Wrap(err, "failed to parse token")
	}

	if !t.Valid {
		return nil, errors.Wrap(err, "invalid token")
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

	if err := g.decode(claims, dest); err != nil {
		return fmt.Errorf("failed to decode payload: %w", err)
	}

	return nil
}

func (g *JWTTokenGenerator) decode(src interface{}, dest interface{}) error {
	tmp, err := json.Marshal(src)
	if err != nil {
		return errors.Wrap(err, "failed to marshal payload")
	}

	if err := json.Unmarshal(tmp, dest); err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}

	return nil
}
