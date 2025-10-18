package services

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var JWT_secret = os.Getenv("JWT_SECRET")

func GetTokenForID(id uint) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	tokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenUnsigned.SignedString([]byte(JWT_secret))
}
