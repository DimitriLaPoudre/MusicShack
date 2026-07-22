package service

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

func (s *AuthService) Login(c context.Context, form model.LoginForm) (string, error) {
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

func (s *AuthService) Logout(c context.Context, token string) error {
	if err := s.session.DeleteSessionByToken(c, token); err != nil {
		return fmt.Errorf("delete session by token %s: %w", token, err)
	}
	return nil
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
