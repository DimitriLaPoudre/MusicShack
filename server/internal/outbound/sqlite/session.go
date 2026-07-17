package sqlite

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite/dto"
	"github.com/google/uuid"
)

func (r *SQLiteRepository) CreateSession(ctx context.Context, s model.Session) (model.Session, error) {
	tx := r.getTx(ctx)

	query := "INSERT INTO sessions (id, user_id, token, expires_at) VALUES (?, ?, ?, ?) RETURNING *"
	query = tx.Rebind(query)

	dbSession := dto.Session{}
	err := tx.GetContext(ctx,
		&dbSession,
		query,
		s.ID, s.UserID, s.Token, s.ExpiresAt,
	)
	if err != nil {
		return model.Session{}, fmt.Errorf("create session: %w", dto.Error(err))
	}

	return dbSession.ToSession(), nil

}

func (r *SQLiteRepository) GetSessionByFilter(ctx context.Context, filter model.SessionFilter) (model.Session, error) {
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
	if filter.Token != nil {
		setParts = append(setParts, "token=?")
		args = append(args, *filter.Token)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM sessions LIMIT 1"
	} else {
		query = fmt.Sprintf("SELECT * FROM sessions WHERE %s LIMIT 1", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbSession := dto.Session{}
	err := tx.GetContext(ctx, &dbSession, query, args...)
	if err != nil {
		return model.Session{}, fmt.Errorf("get session with filter %v: %w", filter, dto.Error(err))
	}

	return dbSession.ToSession(), nil
}

func (r *SQLiteRepository) ListSessionsByFilter(ctx context.Context, filter model.SessionFilter) ([]model.Session, error) {
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
	if filter.Token != nil {
		setParts = append(setParts, "token=?")
		args = append(args, *filter.Token)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM sessions"
	} else {
		query = fmt.Sprintf("SELECT * FROM sessions WHERE %s", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbSessions := []dto.Session{}
	err := tx.SelectContext(ctx, &dbSessions, query, args...)
	if err != nil {
		return []model.Session{}, fmt.Errorf("list session with filter %v: %w", filter, dto.Error(err))
	}

	return dto.SessionsToSessions(dbSessions), nil
}

func (r *SQLiteRepository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM sessions WHERE id = ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("delete session %s: %w", sessionID.String(), dto.Error(err))
	}
	return nil
}

func (r *SQLiteRepository) DeleteSessionByToken(ctx context.Context, token string) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM sessions WHERE token = ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("delete session via token %s: %w", token, dto.Error(err))
	}
	return nil
}

func (r *SQLiteRepository) DeleteSessionExpired(ctx context.Context) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM sessions WHERE expires_at <= ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("delete expired sessions: %w", dto.Error(err))
	}
	return nil
}
