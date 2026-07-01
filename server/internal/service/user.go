package service

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/crypto"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserService struct {
	l    *zerolog.Logger
	cfg  config.LibraryConfig
	repo model.UserRepository
}

func NewUserService(l *zerolog.Logger, cfg config.LibraryConfig, repo model.UserRepository) UserService {
	return UserService{
		l:    l,
		cfg:  cfg,
		repo: repo,
	}
}

func (u *UserService) CreateUser(c context.Context, user *model.User) (*model.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	user.ID = id

	hash, err := crypto.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hash

	//TODO create user space too

	return u.repo.CreateUser(c, user)
}
