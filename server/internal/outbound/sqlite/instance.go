package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite/dto"
	"github.com/google/uuid"
)

func (r *SQLiteRepository) CreateInstance(ctx context.Context, i model.Instance) (model.Instance, error) {
	tx := r.getTx(ctx)

	query := "INSERT INTO instances (id, user_id, provider, plugin, url, ping) VALUES (?, ?, ?, ?, ?, ?) RETURNING *"
	query = tx.Rebind(query)

	dbInstance := dto.Instance{}
	err := tx.GetContext(ctx,
		&dbInstance, query,
		i.ID, i.UserID, i.Provider, i.Plugin, i.Url, i.Ping,
	)
	if err != nil {
		return model.Instance{}, fmt.Errorf("failed to create instance: %w", dto.Error(err))
	}

	return dbInstance.ToInstance(), nil

}

func (r *SQLiteRepository) ListInstancesByFilter(ctx context.Context, filter model.InstanceFilter) ([]model.Instance, error) {
	tx := r.getTx(ctx)

	setParts := []string{}
	args := []any{}

	if filter.ID != nil {
		setParts = append(setParts, "id=?")
		args = append(args, *filter.ID)
	}
	if filter.UserID != nil {
		setParts = append(setParts, "user_id=?")
		args = append(args, *filter.UserID)
	}
	if filter.Provider != nil {
		setParts = append(setParts, "provider=?")
		args = append(args, *filter.Provider)
	}
	if filter.Plugin != nil {
		setParts = append(setParts, "plugin=?")
		args = append(args, *filter.Plugin)
	}
	if filter.Url != nil {
		setParts = append(setParts, "role=?")
		args = append(args, *filter.Url)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM instances ORDER BY ping ASC"
	} else {
		query = fmt.Sprintf("SELECT * FROM instances WHERE %s ORDER BY ping ASC", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbInstances := []dto.Instance{}
	err := tx.SelectContext(ctx, &dbInstances, query, args...)
	if err != nil {
		return []model.Instance{}, fmt.Errorf("failed to create instance: %w", dto.Error(err))
	}

	return dto.InstancesToInstances(dbInstances), nil
}

func (r *SQLiteRepository) DeleteInstance(ctx context.Context, instanceID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM instances WHERE id = ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, instanceID)
	if err != nil {
		return fmt.Errorf("failed to delete instance %s: %w", instanceID.String(), dto.Error(err))
	}

	return nil
}

func (r *SQLiteRepository) DeleteInstanceByUserID(ctx context.Context, instanceID uuid.UUID, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM instances WHERE id = ? AND user_id = ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, instanceID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete instance %s of user %s: %w", instanceID.String(), userID.String(), dto.Error(err))
	}

	return nil
}
