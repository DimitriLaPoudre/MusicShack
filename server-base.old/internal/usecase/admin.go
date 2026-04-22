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
	cfg   *config.AdminConfig
	admin model.AdminRepository
	user  model.UserRepository
}

func NewAdminUseCase(l *zerolog.Logger, cfg *config.AdminConfig, admin model.AdminRepository, user model.UserRepository) AdminUseCase {
	return AdminUseCase{
		l:     l,
		admin: admin,
		user:  user,
	}
}

func (u *AdminUseCase) Login(ctx context.Context, password string) (string, error) {
	admin, err := u.admin.GetAdmin(ctx)
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

func (u *AdminUseCase) ChangeAdminPassword(ctx context.Context, newPassword string, oldPassword string) error {
	return nil
}
