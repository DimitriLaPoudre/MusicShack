package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/crypto"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/token"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	cfg     config.SessionConfig
	user    model.UserRepository
	session model.SessionRepository
}

func NewAuthUseCase(cfg config.SessionConfig, user model.UserRepository, session model.SessionRepository) AuthUseCase {
	return AuthUseCase{
		cfg:     cfg,
		user:    user,
		session: session,
	}
}

func (s *AuthUseCase) Login(c context.Context, form model.LoginForm) (string, error) {
	user, err := s.user.GetUserByFilter(c, model.UserFilter{Username: &form.Username})
	if err != nil {
		return "", fmt.Errorf("get user by username %s: %w", form.Username, err)
	}

	if err := crypto.ComparePassword(user.Password, form.Password); err != nil {
		return "", fmt.Errorf("compare given password with user password hash: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("create id for new session: %w", err)
	}

	token, err := token.GenerateSessionToken()
	if err != nil {
		return "", fmt.Errorf("create token for new session: %w", err)
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

	session, err := s.session.CreateSession(c, newSession)
	if err != nil {
		return "", fmt.Errorf("create new session: %w", err)
	}

	return session.Token, nil

}

func (s *AuthUseCase) Logout(c context.Context, token string) error {
	if err := s.session.DeleteSessionByToken(c, token); err != nil {
		return fmt.Errorf("delete session by token %s: %w", token, err)
	}
	return nil
}
