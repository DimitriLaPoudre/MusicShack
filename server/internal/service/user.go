package service

import (
	"context"
	"fmt"
	"os"

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

func (s *UserService) CreateUser(ctx context.Context, user model.User) (model.User, error) {
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

	if err := s.repo.WithTransaction(ctx, func(ctx context.Context) error {
		user, err = s.repo.CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("create new user: %w", err)
		}

		root, err := os.OpenRoot(s.cfg.Path)
		if err != nil {
			return fmt.Errorf("open download folder: %w", err)
		}

		if err := root.Mkdir(user.ID.String(), 0755); err != nil {
			return fmt.Errorf("create user folder: %w", err)
		}

		return nil
	}); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := s.repo.GetUserByFilter(ctx, model.UserFilter{ID: &id})
	if err != nil {
		return model.User{}, fmt.Errorf("get user %s: %w", id.String(), err)
	}

	return user, nil
}

func (s *UserService) GetUserByFilter(ctx context.Context, filter model.UserFilter) (model.User, error) {
	user, err := s.repo.GetUserByFilter(ctx, filter)
	if err != nil {
		return model.User{}, fmt.Errorf("get user for filter %v: %w", filter, err)
	}

	return user, nil
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.ListUsersByFilter(ctx, model.UserFilter{})
	if err != nil {
		return []model.User{}, fmt.Errorf("list all users: %w", err)
	}

	return users, nil
}

func (s *UserService) ListUsersByFilter(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	users, err := s.repo.ListUsersByFilter(ctx, filter)
	if err != nil {
		return []model.User{}, fmt.Errorf("list users for filter %v: %w", filter, err)
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, partial model.PartialUser) (model.User, error) {
	user, err := s.repo.UpdateUser(ctx, partial)
	if err != nil {
		return model.User{}, fmt.Errorf("update user %s: %w", partial.ID.String(), err)
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user %s: %w", id.String(), err)
	}

	return nil
}
