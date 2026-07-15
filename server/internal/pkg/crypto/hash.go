package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("crypto.HashPassword: bcrypt.GenerateFromPassword: %w", err)
	}
	return string(hash), nil
}

func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("crypto.ComparePassword: bcrypt.CompareHashAndPassword: %w", err)
	}
	return nil
}
