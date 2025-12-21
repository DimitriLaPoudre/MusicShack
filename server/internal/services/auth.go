package services

import (
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func CreateUserSession(userId uint) (*models.UserSession, error) {
	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("CreateUserSession: %w", err)
	}
	expireAt := time.Now().Add(1 * time.Hour)

	session := models.UserSession{
		UserId:    userId,
		Token:     token,
		ExpiresAt: expireAt,
	}

	if err := repository.CreateUserSession(&session); err != nil {
		return nil, fmt.Errorf("CreateUserSession: %w", err)
	}

	return &session, nil
}
