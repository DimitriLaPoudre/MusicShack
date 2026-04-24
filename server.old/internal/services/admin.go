package services

import (
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
)

func CreateAdminSession(expiresIn time.Duration) (string, error) {
	token, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("CreateAdminSession: %w", err)
	}
	expiresAt := time.Now().Add(expiresIn)

	if err := repository.CreateAdminSession(token, expiresAt); err != nil {
		return "", fmt.Errorf("CreateAdminSession: %w", err)
	}

	return token, nil
}

func AdminLogout(c *gin.Context) {
	_ = repository.DeleteAdminSession()
	c.SetCookie("admin_session", "", -1, "/", "", config.HTTPS, true)
}
