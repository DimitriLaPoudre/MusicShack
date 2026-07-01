package usecase

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

type AuthUseCase struct {
	l       *zerolog.Logger
	cfg     config.SessionConfig
	user    model.UserRepository
	session model.SessionRepository
}

func NewAuthUseCase(l *zerolog.Logger, cfg config.SessionConfig, user model.UserRepository, session model.SessionRepository) AuthUseCase {
	return AuthUseCase{
		l:       l,
		cfg:     cfg,
		user:    user,
		session: session,
	}
}

func (s *AuthUseCase) Login(c context.Context, form *model.LoginForm) (string, error) {
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

func (s *AuthUseCase) Logout(c context.Context, token string) error {
	err := s.session.DeleteSessionByToken(c, token)
	return err
}
