package postgres

import (
	"context"
	"fmt"
	"strings"

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

func (r *PostgresRepository) ListInstancesByFilter(ctx context.Context, filter *model.InstanceFilter) ([]*model.Instance, error) {
	if filter == nil {
		return []*model.Instance{}, model.ErrUnknown
	}

	setParts := []string{}
	args := []any{}
	argID := 1

	if filter.ID != nil {
		setParts = append(setParts, fmt.Sprintf("id=$%d", argID))
		args = append(args, *filter.ID)
		argID++
	}
	if filter.UserID != nil {
		setParts = append(setParts, fmt.Sprintf("user_id=$%d", argID))
		args = append(args, *filter.UserID)
		argID++
	}
	if filter.Provider != nil {
		setParts = append(setParts, fmt.Sprintf("provider=$%d", argID))
		args = append(args, *filter.Provider)
		argID++
	}
	if filter.Plugin != nil {
		setParts = append(setParts, fmt.Sprintf("plugin=$%d", argID))
		args = append(args, *filter.Plugin)
		argID++
	}
	if filter.Url != nil {
		setParts = append(setParts, fmt.Sprintf("role=$%d", argID))
		args = append(args, *filter.Url)
		argID++
	}

	var query string
	if argID == 1 {
		query = "SELECT * FROM instances"
	} else {
		query = fmt.Sprintf("SELECT * FROM instances WHERE %s ORDER_BY ping asc", strings.Join(setParts, ", "))
	}

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
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
