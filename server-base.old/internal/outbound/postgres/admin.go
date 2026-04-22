package postgres

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
)

func (r *PostgresRepository) GetAdmin(ctx context.Context) (model.Admin, error) {
	return model.Admin{}, nil
}

func (r *PostgresRepository) UpdateAdminSession(ctx context.Context, token string, expiresAt time.Time) error {
	return nil
}
