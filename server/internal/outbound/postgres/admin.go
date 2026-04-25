package postgres

import (
	"context"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/postgres/dto"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) InitAdmin(ctx context.Context, hashed_password string) (*model.Admin, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, "INSERT INTO admin (id, password, token, expires_at) VALUES ($1, $2, $3, $4) RETURNING *", 1, hashed_password, "", time.Now())
	if err != nil {
		return nil, err
	}

	admin, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Admin])
	if err != nil {
		return nil, err
	}

	return admin.ToAdmin(), nil
}

func (r *PostgresRepository) GetAdmin(ctx context.Context) (*model.Admin, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, "SELECT * FROM admin WHERE id = 1")
	if err != nil {
		return nil, err
	}

	admin, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Admin])
	if err != nil {
		return nil, err
	}

	return admin.ToAdmin(), nil
}

func (r *PostgresRepository) UpdateAdminSession(ctx context.Context, token string, expiresAt time.Time) error {
	tx := r.getTx(ctx)

	_, err := tx.Query(ctx, "UPDATE admin SET token = $1, expires_at = $2 WHERE id = 1", token, expiresAt)
	if err != nil {
		return err
	}

	return nil
}
