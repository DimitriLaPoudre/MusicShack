package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/postgres/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateSession(ctx context.Context, s model.Session) (model.Session, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"INSERT INTO users (id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *",
		s.ID, s.UserID, s.ExpiresAt,
	)
	if err != nil {
		return model.Session{}, err
	}

	dbSession, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return model.Session{}, err
	}

	return dbSession.ToSession(), nil

}

func (r *PostgresRepository) GetSessionByFilter(ctx context.Context, filter model.SessionFilter) (model.Session, error) {
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
	if filter.Token != nil {
		setParts = append(setParts, fmt.Sprintf("token=$%d", argID))
		args = append(args, *filter.Token)
		argID++
	}

	var query string
	if argID == 1 {
		query = "SELECT * FROM sessions LIMIT 1"
	} else {
		query = fmt.Sprintf("SELECT * FROM sessions WHERE %s LIMIT 1", strings.Join(setParts, ", "))
	}

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return model.Session{}, err
	}

	dbSession, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return model.Session{}, err
	}

	return dbSession.ToSession(), nil
}

func (r *PostgresRepository) ListSessionsByFilter(ctx context.Context, filter model.SessionFilter) ([]model.Session, error) {
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
	if filter.Token != nil {
		setParts = append(setParts, fmt.Sprintf("token=$%d", argID))
		args = append(args, *filter.Token)
		argID++
	}

	var query string
	if argID == 1 {
		query = "SELECT * FROM sessions"
	} else {
		query = fmt.Sprintf("SELECT * FROM sessions WHERE %s", strings.Join(setParts, ", "))
	}

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return []model.Session{}, err
	}

	dbSessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return []model.Session{}, err
	}

	return dto.SessionsToSessions(dbSessions), nil
}

func (r *PostgresRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM sessions WHERE id = $1", sessionID); err != nil {
		return dto.Error(err)
	}
	return nil
}

func (r *PostgresRepository) DeleteSessionByToken(ctx context.Context, token string) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM sessions WHERE token = $1", token); err != nil {
		return dto.Error(err)
	}
	return nil
}

func (r *PostgresRepository) DeleteSessionExpired(ctx context.Context) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM sessions WHERE expires_at <= $1", time.Now()); err != nil {
		return dto.Error(err)
	}
	return nil
}
