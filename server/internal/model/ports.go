package model

import (
	"context"
	"time"
)

type AdminRepository interface {
	InitAdmin(ctx context.Context, hashed_password string) (*Admin, error)
	GetAdmin(ctx context.Context) (*Admin, error)
	UpdateAdminSession(ctx context.Context, token string, expiresAt time.Time) error
}

type UserRepository interface {
}
