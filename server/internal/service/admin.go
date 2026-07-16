package service

import (
	"context"
	"errors"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
)

func InitAdmin(c context.Context, cfg config.AdminConfig, user *UserService, repo model.UserRepository) error {
	role := model.UserRoleAdmin
	if _, err := repo.GetUserByFilter(c, model.UserFilter{Role: &role}); err == nil || !errors.Is(err, model.ErrNotFound) {
		return err
	}

	if _, err := user.CreateUser(c, model.User{
		Username: cfg.DefaultUsername,
		Password: cfg.DefaultPassword,
		Role:     model.UserRoleAdmin,
	}); err != nil {
		return err
	}

	return nil
}
