package service

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/rs/zerolog"
)

type AuthService struct {
	l       *zerolog.Logger
	cfg     config.SessionConfig
	user    model.UserRepository
	session model.SessionRepository
}

func NewAuthService(l *zerolog.Logger, cfg config.SessionConfig, user model.UserRepository, session model.SessionRepository) AuthService {
	return AuthService{
		l:       l,
		cfg:     cfg,
		user:    user,
		session: session,
	}
}

func (s *AuthService) Authenticate(c context.Context, tkn string) (model.User, error) {
	session, err := s.session.GetSessionByToken(c, tkn)
	if err != nil {
		return model.User{}, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return model.User{}, model.ErrBadToken
	}

	user, err := s.user.GetUserByFilter(c, model.UserFilter{ID: &session.UserID})
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *AuthService) IsAdmin(c context.Context, user model.User) bool {
	return user.Role == model.UserRoleAdmin
}

func (s *AuthService) IsUser(c context.Context, user model.User) bool {
	return user.Role == model.UserRoleUser
}
