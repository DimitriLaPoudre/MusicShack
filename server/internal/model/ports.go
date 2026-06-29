package model

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	ListAllUsers(ctx context.Context) ([]*User, error)
	UpdateUser(ctx context.Context, partialUser *PartialUser) (*User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, s *Session) (*Session, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*Session, error)
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
}
