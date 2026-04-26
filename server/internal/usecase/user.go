package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/crypto"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserUseCase struct {
	l    *zerolog.Logger
	cfg  config.LibraryConfig
	repo model.UserRepository
}

func NewUserUseCase(l *zerolog.Logger, cfg config.LibraryConfig, repo model.UserRepository) UserUseCase {
	return UserUseCase{
		l:    l,
		cfg:  cfg,
		repo: repo,
	}
}

func (u *UserUseCase) CreateUser(c context.Context, user *model.User) (*model.User, error) {
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

func (u *UserUseCase) GetUserByID(c context.Context, id uuid.UUID) (*model.User, error) {
	return u.repo.GetUserByID(c, id)
}

func (u *UserUseCase) ListAllUsers(c context.Context) ([]*model.User, error) {
	return u.repo.ListAllUsers(c)
}

func (u *UserUseCase) UpdateUser(c context.Context, user *model.PartialUser) (*model.User, error) {
	return u.repo.UpdateUser(c, user)
}

func (u *UserUseCase) DeleteUser(c context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(c, id)
}
