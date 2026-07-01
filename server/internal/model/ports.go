package model

import (
	"context"
	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserWithFilter(ctx context.Context, filter *FilterUser) (*User, error)
	ListUsersWithFilter(ctx context.Context, filter *FilterUser) ([]*User, error)
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
	DeleteSessionByToken(ctx context.Context, token string) error
}

type InstanceRepository interface {
	CreateInstance(ctx context.Context, i *Instance) (*Instance, error)
	GetInstance(ctx context.Context, id uuid.UUID) (*Instance, error)
	ListInstancesByUserID(ctx context.Context, userID uuid.UUID) ([]*Instance, error)
	DeleteInstance(ctx context.Context, id uuid.UUID) error
	DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}
