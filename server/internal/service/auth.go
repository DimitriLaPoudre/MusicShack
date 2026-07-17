package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
)

type AuthService struct {
	cfg     config.SessionConfig
	user    model.UserRepository
	session model.SessionRepository
}

func NewAuthService(cfg config.SessionConfig, user model.UserRepository, session model.SessionRepository) AuthService {
	return AuthService{
		cfg:     cfg,
		user:    user,
		session: session,
	}
}

func (s *AuthService) Authenticate(c context.Context, tkn string) (model.User, error) {
	session, err := s.session.GetSessionByFilter(c, model.SessionFilter{Token: &tkn})
	if err != nil {
		return model.User{}, fmt.Errorf("get token for authentication: %w", err)
	}

	if session.ExpiresAt.Before(time.Now()) {
		return model.User{}, model.ErrExpiredToken
	}

	user, err := s.user.GetUserByFilter(c, model.UserFilter{ID: &session.UserID})
	if err != nil {
		return model.User{}, fmt.Errorf("get user %s from session %s: %w", session.UserID, session.ID, err)
	}

	return user, nil
}

func (s *AuthService) IsAdmin(c context.Context, user model.User) bool {
	return user.Role == model.UserRoleAdmin
}

func (s *AuthService) IsUser(c context.Context, user model.User) bool {
	return user.Role == model.UserRoleUser
}
