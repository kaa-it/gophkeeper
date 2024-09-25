package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/kaa-it/gophkeeper/internal/server/domain"
)

// ErrInvalidTokenClaims indicates that the JWT token claims are invalid or cannot be parsed.
var ErrInvalidTokenClaims = errors.New("invalid token claims")

// JWTManager is responsible for generating and validating JWT tokens.
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// UserClaims defines the structure for JWT claims including standard claims and username.
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

// NewJWTManager creates a new instance of JWTManager with the provided secret key and token duration.
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Generate creates a JWT token for the specified user and returns it as a string.
// It returns an error if the token generation fails.
func (m *JWTManager) Generate(user *domain.User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		Username: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Verify validates the provided JWT token and returns the user
// claims if the token is valid; otherwise, returns an error.
func (m *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(_ *jwt.Token) (any, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
