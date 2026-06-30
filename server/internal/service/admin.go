package service

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
)

func InitAdmin(c context.Context, cfg config.AdminConfig, user UserService) error {
	role := model.UserRoleAdmin
	if _, err := user.GetUserWithFilter(c, &model.FilterUser{Role: &role}); err != nil {
		return err
	}

	if _, err := user.CreateUser(c, &model.User{
		Username: cfg.DefaultUsername,
		Password: cfg.DefaultPassword,
		Role:     model.UserRoleAdmin,
	}); err != nil {
		return err
	}

	return nil
}
