package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("session token generation: %v", err)
	}
	return hex.EncodeToString(b), nil
}
