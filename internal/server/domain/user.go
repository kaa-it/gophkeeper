package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string
	HashedPassword string
}

func NewUser(username, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	return &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
	}, nil
}

func (u *User) IsCorrectPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)) == nil
}

func (u *User) Clone() *User {
	return &User{
		Username:       u.Username,
		HashedPassword: u.HashedPassword,
	}
}
