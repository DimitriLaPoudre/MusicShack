package usecase

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/rs/zerolog"
)

type SessionUseCase struct {
	l       *zerolog.Logger
	cfg     config.AdminConfig
	user    model.UserRepository
	session model.SessionRepository
}

func NewSessionUseCase(l *zerolog.Logger, cfg config.AdminConfig, user model.UserRepository, session model.SessionRepository) SessionUseCase {
	return SessionUseCase{
		l:       l,
		cfg:     cfg,
		user:    user,
		session: session,
	}
}

func (s *SessionUseCase) Authenticate(c context.Context, tkn string) (*model.User, error) {
	session, err := s.session.GetSessionByToken(c, tkn)
	if err != nil {
		return nil, err
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, model.ErrBadToken
	}

	user, err := s.user.GetUserByID(c, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SessionUseCase) IsAdmin(c context.Context, user *model.User) bool {
	return user.Role == model.UserRoleAdmin
}

func (s *SessionUseCase) IsUser(c context.Context, user *model.User) bool {
	return user.Role == model.UserRoleUser
}
