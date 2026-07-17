package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
)

func InitAdmin(c context.Context, cfg config.AdminConfig, user *UserService, repo model.UserRepository) error {
	role := model.UserRoleAdmin
	if _, err := repo.GetUserByFilter(c, model.UserFilter{Role: &role}); !errors.Is(err, model.ErrNotFound) {
		if err != nil {
			return fmt.Errorf("check if admin already exist: %w", err)
		}
		return model.ErrAdminAlreadyExist
	}

	if _, err := user.CreateUser(c, model.User{
		Username: cfg.DefaultUsername,
		Password: cfg.DefaultPassword,
		Role:     model.UserRoleAdmin,
	}); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}
