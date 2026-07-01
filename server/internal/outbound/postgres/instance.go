package postgres

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/postgres/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateInstance(ctx context.Context, i *model.Instance) (*model.Instance, error) {
	if i == nil {
		return nil, model.ErrInvalidInput
	}
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"INSERT INTO instances (id, user_id, provider, plugin, url, ping) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *",
		i.ID, i.UserID, i.Provider, i.Plugin, i.Url, i.Ping,
	)
	if err != nil {
		return nil, err
	}

	dbInstance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Instance])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbInstance.ToInstance(), nil

}

func (r *PostgresRepository) GetInstance(ctx context.Context, id uuid.UUID) (*model.Instance, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM instances WHERE id = $1 LIMIT 1",
		id)
	if err != nil {
		return nil, err
	}

	dbInstance, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Instance])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbInstance.ToInstance(), nil
}

func (r *PostgresRepository) ListInstancesByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Instance, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM instances WHERE user_id = $1 ORDER_BY ping asc",
		userID)
	if err != nil {
		return nil, err
	}

	dbInstances, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.Instance])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dto.InstancesToInstances(dbInstances), nil
}

func (r *PostgresRepository) DeleteInstance(ctx context.Context, id uuid.UUID) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM instances WHERE id = $1", id); err != nil {
		return dto.Error(err)
	}
	return nil
}

func (r *PostgresRepository) DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM instances WHERE id = $1 AND user_id = $2", id, userID); err != nil {
		return dto.Error(err)
	}
	return nil
}
