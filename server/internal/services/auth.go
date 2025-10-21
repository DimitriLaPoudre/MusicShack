package services

import (
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetTokenForID(id uint) (string, error) {
	claims := models.JwtClaims{
		Id: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenUnsigned := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return tokenUnsigned.SignedString(config.JWT_SECRET)
}

func GetTokenContent(token string) (*models.JwtClaims, error) {
	tokenUnsigned, err := jwt.ParseWithClaims(token, &models.JwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method not valid: %v", token.Header["alg"])
		}
		return []byte(config.JWT_SECRET), nil
	})
	if err != nil {
		return nil, err
	}
	if !tokenUnsigned.Valid {
		return nil, fmt.Errorf("token invalid")
	}
	claims, ok := tokenUnsigned.Claims.(*models.JwtClaims)
	if !ok {
		return nil, fmt.Errorf("token content invalid")
	}
	return claims, nil
}
