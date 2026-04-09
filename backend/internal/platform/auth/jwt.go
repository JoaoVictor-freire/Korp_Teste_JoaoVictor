package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenSigner struct {
	secret []byte
}

func NewTokenSignerFromEnv() TokenSigner {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Dev default. For production, always set JWT_SECRET.
		secret = "dev-secret-change-me"
	}

	return TokenSigner{secret: []byte(secret)}
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (s TokenSigner) Sign(userID string, email string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s TokenSigner) Verify(rawToken string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(rawToken, &Claims{}, func(_ *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
