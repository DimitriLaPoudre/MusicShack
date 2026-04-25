package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AdminRepository interface {
	InitAdmin(ctx context.Context, hashed_password string) (*Admin, error)
	GetAdmin(ctx context.Context) (*Admin, error)
	UpdateAdminSession(ctx context.Context, token string, expiresAt time.Time) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	ListAllUsers(ctx context.Context) ([]*User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
