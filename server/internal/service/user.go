package service

import (
	"context"
	"fmt"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/crypto"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type UserService struct {
	cfg  config.DownloadConfig
	repo model.UserRepository
}

func NewUserService(cfg config.DownloadConfig, repo model.UserRepository) UserService {
	return UserService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *UserService) CreateUser(c context.Context, user model.User) (model.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return model.User{}, fmt.Errorf("create id for new user: %w", err)
	}
	user.ID = id

	hash, err := crypto.HashPassword(user.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("create password hash for new user: %w", err)
	}
	user.Password = hash

	if err := user.Role.IsValid(); err != nil {
		return model.User{}, fmt.Errorf("validate role of new user: %w", err)
	}

	//TODO create user space too

	user, err = s.repo.CreateUser(c, user)
	if err != nil {
		return model.User{}, fmt.Errorf("create new user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(c context.Context, id uuid.UUID) (model.User, error) {
	user, err := s.repo.GetUserByFilter(c, model.UserFilter{ID: &id})
	if err != nil {
		return model.User{}, fmt.Errorf("get user %s: %w", id.String(), err)
	}

	return user, nil
}

func (s *UserService) GetUserByFilter(c context.Context, filter model.UserFilter) (model.User, error) {
	user, err := s.repo.GetUserByFilter(c, filter)
	if err != nil {
		return model.User{}, fmt.Errorf("get user for filter %v: %w", filter, err)
	}

	return user, nil
}

func (s *UserService) ListAllUsers(c context.Context) ([]model.User, error) {
	users, err := s.repo.ListUsersByFilter(c, model.UserFilter{})
	if err != nil {
		return []model.User{}, fmt.Errorf("list all users: %w", err)
	}

	return users, nil
}

func (s *UserService) ListUsersWithFilter(c context.Context, filter model.UserFilter) ([]model.User, error) {
	users, err := s.repo.ListUsersByFilter(c, filter)
	if err != nil {
		return []model.User{}, fmt.Errorf("list users for filter %v: %w", filter, err)
	}

	return users, nil
}

func (s *UserService) UpdateUser(c context.Context, partial model.PartialUser) (model.User, error) {
	user, err := s.repo.UpdateUser(c, partial)
	if err != nil {
		return model.User{}, fmt.Errorf("update user %s: %w", partial.ID.String(), err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(c context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteUser(c, id); err != nil {
		return fmt.Errorf("delete user %s: %w", id.String(), err)
	}

	return nil
}
