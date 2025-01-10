package auth

import "github.com/golang-jwt/jwt/v5"

// Authenticator - this kinda defeats the purpose of the abstraction, it is tied to JWT
type Authenticator interface {
	GenerateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}
