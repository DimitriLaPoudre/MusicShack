package utils

import (
	"fmt"
	"path/filepath"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
)

func GetUserPath(userId uint) (string, error) {
	if user, err := repository.GetUserByID(userId); err != nil {
		return "", fmt.Errorf("utils.GetUserPath: %w", err)
	} else {
		return filepath.Join(config.LIBRARY_PATH, user.Username), nil
	}
}
