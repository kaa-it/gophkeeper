package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user with a username, login, and hashed password.
type User struct {
	Username       string
	Login          string
	HashedPassword string
}

// NewUser creates a new User instance with the given username, login, and password, hashing the password.
func NewUser(username, login, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	return &User{
		Username:       username,
		Login:          login,
		HashedPassword: string(hashedPassword),
	}, nil
}

// IsCorrectPassword verifies if the provided plaintext password matches the stored hashed password.
func (u *User) IsCorrectPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)) == nil
}

// Clone creates and returns a deep copy of the User instance.
func (u *User) Clone() *User {
	return &User{
		Username:       u.Username,
		Login:          u.Login,
		HashedPassword: u.HashedPassword,
	}
}
