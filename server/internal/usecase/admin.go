package usecase

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/crypto"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/token"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/rs/zerolog"
)

type AdminUseCase struct {
	l     *zerolog.Logger
	cfg   config.AdminConfig
	admin model.AdminRepository
}

func NewAdminUseCase(l *zerolog.Logger, cfg config.AdminConfig, admin model.AdminRepository) AdminUseCase {
	return AdminUseCase{
		l:     l,
		cfg:   cfg,
		admin: admin,
	}
}

func (u *AdminUseCase) Login(ctx context.Context, password string) (string, error) {
	admin, err := u.GetAdmin(ctx)
	if err != nil {
		return "", err
	}

	if err := crypto.ComparePassword(admin.Password, password); err != nil {
		return "", err
	}

	tkn, err := token.GenerateSessionToken()
	if err != nil {
		return "", err
	}

	if err := u.admin.UpdateAdminSession(ctx, tkn, time.Now().Add(u.cfg.TokenDuration)); err != nil {
		return "", err
	}

	return tkn, nil
}

func (u *AdminUseCase) Authenticate(ctx context.Context, tkn string) error {
	admin, err := u.admin.GetAdmin(ctx)
	if err != nil {
		return err
	}

	if tkn != admin.Token || admin.ExpiresAt.Before(time.Now()) {
		return err
	}
	return nil
}

func (u *AdminUseCase) ChangeAdminPassword(ctx context.Context, newPassword string, oldPassword string) error {
	return nil
}

func (u *AdminUseCase) GetAdmin(ctx context.Context) (*model.Admin, error) {
	admin, err := u.admin.GetAdmin(ctx)
	if err != nil {
		var hash string
		hash, err = crypto.HashPassword(u.cfg.DefaultPassword)
		if err != nil {
			return nil, err
		}
		admin, err = u.admin.InitAdmin(ctx, hash)
		if err != nil {
			return nil, err
		}
	}
	return admin, err
}
