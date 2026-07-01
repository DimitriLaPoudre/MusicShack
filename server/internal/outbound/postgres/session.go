package postgres

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/postgres/dto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *PostgresRepository) CreateSession(ctx context.Context, s *model.Session) (*model.Session, error) {
	if s == nil {
		return nil, model.ErrUnknown
	}
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"INSERT INTO users (id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING *",
		s.ID, s.UserID, s.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	dbSession, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbSession.ToSession(), nil

}

func (r *PostgresRepository) GetSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM sessions WHERE id = $1 LIMIT 1",
		sessionID)
	if err != nil {
		return nil, err
	}

	dbSession, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbSession.ToSession(), nil
}

func (r *PostgresRepository) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM sessions WHERE token = $1",
		token)
	if err != nil {
		return nil, err
	}

	dbSession, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbSession.ToSession(), nil
}

func (r *PostgresRepository) GetSessionByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM sessions WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}

	dbSessions, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.Session])
	if err := dto.Error(err); err != nil {
		return nil, err
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
