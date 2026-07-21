package usecase

import (
	"context"
	"fmt"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/google/uuid"
)

type UserUseCase struct {
	user *service.UserService
	repo model.UserRepository
}

func NewUserUseCase(user *service.UserService, repo model.UserRepository) UserUseCase {
	return UserUseCase{
		user: user,
		repo: repo,
	}
}

func (u *UserUseCase) CreateUser(c context.Context, user model.User) (model.User, error) {
	return u.user.CreateUser(c, user)
}

func (u *UserUseCase) GetUserByID(c context.Context, id uuid.UUID) (model.User, error) {
	user, err := u.repo.GetUserByFilter(c, model.UserFilter{ID: &id})
	if err != nil {
		return model.User{}, fmt.Errorf("get user %s: %w", id.String(), err)
	}

	return user, nil
}

func (u *UserUseCase) GetUserByFilter(c context.Context, filter model.UserFilter) (model.User, error) {
	user, err := u.repo.GetUserByFilter(c, filter)
	if err != nil {
		return model.User{}, fmt.Errorf("get user for filter %v: %w", filter, err)
	}

	return user, nil
}

func (u *UserUseCase) ListAllUsers(c context.Context) ([]model.User, error) {
	users, err := u.repo.ListUsersByFilter(c, model.UserFilter{})
	if err != nil {
		return []model.User{}, fmt.Errorf("list all users: %w", err)
	}

	return users, nil
}

func (u *UserUseCase) ListUsersWithFilter(c context.Context, filter model.UserFilter) ([]model.User, error) {
	users, err := u.repo.ListUsersByFilter(c, filter)
	if err != nil {
		return []model.User{}, fmt.Errorf("list users for filter %v: %w", filter, err)
	}

	return users, nil
}

func (u *UserUseCase) UpdateUser(c context.Context, partial model.PartialUser) (model.User, error) {
	user, err := u.repo.UpdateUser(c, partial)
	if err != nil {
		return model.User{}, fmt.Errorf("update user %s: %w", partial.ID.String(), err)
	}

	return user, nil
}

func (u *UserUseCase) DeleteUser(c context.Context, id uuid.UUID) error {
	if err := u.repo.DeleteUser(c, id); err != nil {
		return fmt.Errorf("delete user %s: %w", id.String(), err)
	}

	return nil
}
