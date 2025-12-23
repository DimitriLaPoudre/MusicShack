package repository

import (
	"fmt"
	"time"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func CreateUserSession(userSession *models.UserSession) error {
	if err := database.DB.Create(userSession).Error; err != nil {
		return fmt.Errorf("repository.CreateUserSession: %w", err)
	}
	return nil
}

func FindUserSessionByToken(token string) (*models.UserSession, error) {
	var userSession models.UserSession
	if err := database.DB.First(&userSession, "token = ? AND expires_at > ?", token, time.Now()).Error; err != nil {
		return nil, fmt.Errorf("repository.FindUserSession: %w", err)
	}
	return &userSession, nil
}

func DeleteUserSession(userId uint, token string) error {
	if err := database.DB.Delete(&models.UserSession{}, "user_id = ? AND token = ?", userId, token).Error; err != nil {
		return fmt.Errorf("repository.DeleteUserSession: %w", err)
	}
	return nil
}

func DeleteExpiredUserSession() error {
	if err := database.DB.Delete(&models.UserSession{}, "expires_at <= ?", time.Now()).Error; err != nil {
		return fmt.Errorf("repository.DeleteExpiredUserSession: %w", err)
	}
	return nil
}
