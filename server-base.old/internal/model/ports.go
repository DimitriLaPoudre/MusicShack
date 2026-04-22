package model

import (
	"context"
	"time"
)

type AdminRepository interface {
	GetAdmin(ctx context.Context) (Admin, error)
	UpdateAdminSession(ctx context.Context, token string, expiresAt time.Time) error
}

type UserRepository interface {
}
