package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type MockAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-iss",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
}

func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{}
}

func (a *MockAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	return token.SignedString([]byte(secret))
}

func (a *MockAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
