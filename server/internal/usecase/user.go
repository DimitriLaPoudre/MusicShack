package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type UserUseCase struct {
	l    *zerolog.Logger
	cfg  config.LibraryConfig
	user *service.UserService
}

func NewUserUseCase(l *zerolog.Logger, cfg config.LibraryConfig, user *service.UserService) UserUseCase {
	return UserUseCase{
		l:    l,
		cfg:  cfg,
		user: user,
	}
}

func (u *UserUseCase) CreateUser(c context.Context, user *model.User) (*model.User, error) {
	return u.user.CreateUser(c, user)
}

func (u *UserUseCase) GetUserByID(c context.Context, id uuid.UUID) (*model.User, error) {
	return u.user.GetUserByID(c, id)
}

func (u *UserUseCase) GetUserWithFilter(c context.Context, filter *model.FilterUser) (*model.User, error) {
	return u.user.GetUserWithFilter(c, filter)
}

func (u *UserUseCase) ListUsersWithFilter(c context.Context, filter *model.FilterUser) ([]*model.User, error) {
	return u.user.ListUsersWithFilter(c, filter)
}

func (u *UserUseCase) ListAllUsers(c context.Context) ([]*model.User, error) {
	return u.user.ListAllUsers(c)
}

func (u *UserUseCase) UpdateUser(c context.Context, user *model.PartialUser) (*model.User, error) {
	return u.user.UpdateUser(c, user)
}

func (u *UserUseCase) DeleteUser(c context.Context, id uuid.UUID) error {
	return u.user.DeleteUser(c, id)
}
