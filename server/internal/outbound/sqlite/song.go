package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite/dto"
	"github.com/google/uuid"
)

func (r *SQLiteRepository) CreateSong(ctx context.Context, s model.Song) (model.Song, error) {
	tx := r.getTx(ctx)

	query := "INSERT INTO songs (id, user_id, isrc, path, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING *"
	query = tx.Rebind(query)

	dbSong := dto.Song{}
	err := tx.GetContext(ctx,
		&dbSong,
		query,
		s.ID, s.UserID, s.ISRC, s.Path, s.UpdatedAt,
	)
	if err != nil {
		return model.Song{}, dto.Error(err)
	}

	return dbSong.ToSong(), nil
}

func (r *SQLiteRepository) GetSongByFilter(ctx context.Context, filter model.SongFilter) (model.Song, error) {
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

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM songs LIMIT 1"
	} else {
		query = fmt.Sprintf("SELECT * FROM songs WHERE %s LIMIT 1", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbSong := dto.Song{}
	err := tx.GetContext(ctx, &dbSong, query, args...)
	if err != nil {
		return model.Song{}, dto.Error(err)
	}

	return dbSong.ToSong(), nil
}

func (r *SQLiteRepository) ListSongsByFilter(ctx context.Context, filter model.SongFilter) ([]model.Song, error) {
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

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM songs"
	} else {
		query = fmt.Sprintf("SELECT * FROM songs WHERE %s", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbSongs := []dto.Song{}
	err := tx.SelectContext(ctx, &dbSongs, query, args...)
	if err != nil {
		return []model.Song{}, dto.Error(err)
	}

	return dto.SongsToSongs(dbSongs), nil
}

func (r *SQLiteRepository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM songs WHERE id = ?"
	query = tx.Rebind(query)

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return dto.Error(err)
	}
	if n, err := result.RowsAffected(); err != nil {
		return dto.Error(err)
	} else if n == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *SQLiteRepository) DeleteSongByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM songs WHERE id = ? AND user_id = ?"
	query = tx.Rebind(query)

	result, err := tx.ExecContext(ctx, query, id, userID)
	if err != nil {
		return dto.Error(err)
	}
	if n, err := result.RowsAffected(); err != nil {
		return dto.Error(err)
	} else if n == 0 {
		return model.ErrNotFound
	}

	return nil
}
