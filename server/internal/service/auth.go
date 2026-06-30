package service

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/crypto"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/token"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
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

func (s *AuthService) Authenticate(c context.Context, tkn string) (*model.User, error) {
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

func (s *AuthService) IsAdmin(c context.Context, user *model.User) bool {
	return user.Role == model.UserRoleAdmin
}

func (s *AuthService) IsUser(c context.Context, user *model.User) bool {
	return user.Role == model.UserRoleUser
}

func (s *AuthService) Login(c context.Context, form *model.LoginForm) (string, error) {
	user, err := s.user.GetUserByUsername(c, form.Username)
	if err != nil {
		return "", err
	}

	if err := crypto.ComparePassword(user.Password, form.Password); err != nil {
		return "", err
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	token, err := token.GenerateSessionToken()
	if err != nil {
		return "", err
	}

	var expiresAt time.Time
	if form.Remember {
		expiresAt = time.Now().Add(s.cfg.RememberExp)
	} else {
		expiresAt = time.Now().Add(s.cfg.Exp)
	}

	newSession := model.Session{
		ID:        id,
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	session, err := s.session.CreateSession(c, &newSession)
	if err != nil {
		return "", err
	}

	return session.Token, nil

}

func (s *AuthService) Logout(c context.Context, token string) error {
	err := s.session.DeleteSessionByToken(c, token)
	return err
}
