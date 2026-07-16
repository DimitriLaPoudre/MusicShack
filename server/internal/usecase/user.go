package usecase

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type UserUseCase struct {
	cfg  config.LibraryConfig
	user *service.UserService
	repo model.UserRepository
}

func NewUserUseCase(cfg config.LibraryConfig, user *service.UserService, repo model.UserRepository) UserUseCase {
	return UserUseCase{
		cfg:  cfg,
		user: user,
		repo: repo,
	}
}

func (u *UserUseCase) CreateUser(c context.Context, user model.User) (model.User, error) {
	return u.user.CreateUser(c, user)
}

func (u *UserUseCase) GetUserByID(c context.Context, id uuid.UUID) (model.User, error) {
	return u.repo.GetUserByFilter(c, model.UserFilter{ID: &id})
}

func (u *UserUseCase) GetUserByFilter(c context.Context, filter model.UserFilter) (model.User, error) {
	return u.repo.GetUserByFilter(c, filter)
}

func (u *UserUseCase) ListAllUsers(c context.Context) ([]model.User, error) {
	return u.repo.ListUsersByFilter(c, model.UserFilter{})
}

func (u *UserUseCase) ListUsersWithFilter(c context.Context, filter model.UserFilter) ([]model.User, error) {
	return u.repo.ListUsersByFilter(c, filter)
}

func (u *UserUseCase) UpdateUser(c context.Context, user model.PartialUser) (model.User, error) {
	return u.repo.UpdateUser(c, user)
}

func (u *UserUseCase) DeleteUser(c context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(c, id)
}
