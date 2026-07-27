package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite/dto"
	"github.com/google/uuid"
)

func (r *SQLiteRepository) CreateFollow(ctx context.Context, f model.Follow) (model.Follow, error) {
	tx := r.getTx(ctx)

	query := "INSERT INTO follows (id, user_id, provider, artist_id, artist_name, artist_picture_url, featuring) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *"
	query = tx.Rebind(query)

	dbFollow := dto.Follow{}
	err := tx.GetContext(ctx,
		&dbFollow, query,
		f.ID, f.UserID, f.Provider, f.ArtistID, f.ArtistName, f.ArtistPictureURL, f.Featuring,
	)
	if err != nil {
		return model.Follow{}, dto.Error(err)
	}

	return dbFollow.ToFollow(), nil

}

func (r *SQLiteRepository) GetFollowByFilter(ctx context.Context, filter model.FollowFilter) (model.Follow, error) {
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
	if filter.ArtistID != nil {
		setParts = append(setParts, "artist_id=?")
		args = append(args, *filter.ArtistID)
	}
	if filter.ArtistName != nil {
		setParts = append(setParts, "artist_name=?")
		args = append(args, *filter.ArtistName)
	}
	if filter.ArtistPictureURL != nil {
		setParts = append(setParts, "artist_picture_url=?")
		args = append(args, *filter.ArtistPictureURL)
	}
	if filter.Featuring != nil {
		setParts = append(setParts, "featuring=?")
		args = append(args, *filter.Featuring)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM follows LIMIT 1"
	} else {
		query = fmt.Sprintf("SELECT * FROM follows WHERE %s LIMIT 1", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbFollow := dto.Follow{}
	err := tx.GetContext(ctx, &dbFollow, query, args...)
	if err != nil {
		return model.Follow{}, dto.Error(err)
	}

	return dbFollow.ToFollow(), nil
}

func (r *SQLiteRepository) ListFollowsByFilter(ctx context.Context, filter model.FollowFilter) ([]model.Follow, error) {
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
	if filter.ArtistID != nil {
		setParts = append(setParts, "artist_id=?")
		args = append(args, *filter.ArtistID)
	}
	if filter.ArtistName != nil {
		setParts = append(setParts, "artist_name=?")
		args = append(args, *filter.ArtistName)
	}
	if filter.ArtistPictureURL != nil {
		setParts = append(setParts, "artist_picture_url=?")
		args = append(args, *filter.ArtistPictureURL)
	}
	if filter.Featuring != nil {
		setParts = append(setParts, "featuring=?")
		args = append(args, *filter.Featuring)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM follows"
	} else {
		query = fmt.Sprintf("SELECT * FROM follows WHERE %s", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbFollows := []dto.Follow{}
	err := tx.SelectContext(ctx, &dbFollows, query, args...)
	if err != nil {
		return []model.Follow{}, dto.Error(err)
	}

	return dto.FollowsToFollows(dbFollows), nil
}

func (r *SQLiteRepository) DeleteFollow(ctx context.Context, followID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM follows WHERE id = ?"
	query = tx.Rebind(query)

	result, err := tx.ExecContext(ctx, query, followID)
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

func (r *SQLiteRepository) DeleteFollowByUserID(ctx context.Context, followID uuid.UUID, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM follows WHERE id = ? AND user_id = ?"
	query = tx.Rebind(query)

	result, err := tx.ExecContext(ctx, query, followID, userID)
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
