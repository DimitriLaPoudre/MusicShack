package services

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func ConstructUser(user *models.RequestUser) *models.User {
	return &models.User{
		Username: user.Username,
		Password: user.Password,
	}
}
