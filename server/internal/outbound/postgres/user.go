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

func (r *PostgresRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"INSERT INTO users (id, username, password, hi_res, role) VALUES ($1, $2, $3, $4, $5) RETURNING *",
		user.ID, user.Username, user.Password, user.HiRes, user.Role)
	if err != nil {
		return model.User{}, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return model.User{}, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) GetUserByFilter(ctx context.Context, filter model.FilterUser) (model.User, error) {
	setParts := []string{}
	args := []any{}
	argID := 1

	if filter.ID != nil {
		setParts = append(setParts, fmt.Sprintf("id=$%d", argID))
		args = append(args, *filter.ID)
		argID++
	}
	if filter.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username=$%d", argID))
		args = append(args, *filter.Username)
		argID++
	}
	if filter.Password != nil {
		setParts = append(setParts, fmt.Sprintf("password=$%d", argID))
		args = append(args, *filter.Password)
		argID++
	}
	if filter.HiRes != nil {
		setParts = append(setParts, fmt.Sprintf("hi_res=$%d", argID))
		args = append(args, *filter.HiRes)
		argID++
	}
	if filter.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role=$%d", argID))
		args = append(args, *filter.Role)
		argID++
	}

	var query string
	if argID == 1 {
		query = "SELECT * FROM users"
	} else {
		query = fmt.Sprintf("SELECT * FROM users WHERE %s LIMIT 1", strings.Join(setParts, ", "))
	}

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return model.User{}, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return model.User{}, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) ListUsersByFilter(ctx context.Context, filter model.FilterUser) ([]model.User, error) {
	setParts := []string{}
	args := []any{}
	argID := 1

	if filter.ID != nil {
		setParts = append(setParts, fmt.Sprintf("id=$%d", argID))
		args = append(args, *filter.ID)
		argID++
	}
	if filter.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username=$%d", argID))
		args = append(args, *filter.Username)
		argID++
	}
	if filter.Password != nil {
		setParts = append(setParts, fmt.Sprintf("password=$%d", argID))
		args = append(args, *filter.Password)
		argID++
	}
	if filter.HiRes != nil {
		setParts = append(setParts, fmt.Sprintf("hi_res=$%d", argID))
		args = append(args, *filter.HiRes)
		argID++
	}
	if filter.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role=$%d", argID))
		args = append(args, *filter.Role)
		argID++
	}

	var query string
	if argID == 1 {
		query = "SELECT * FROM users"
	} else {
		query = fmt.Sprintf("SELECT * FROM users WHERE %s", strings.Join(setParts, ", "))
	}

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return []model.User{}, err
	}

	dbUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return []model.User{}, err
	}

	return dto.UsersToUsers(dbUsers), nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, partialUser model.PartialUser) (model.User, error) {
	setParts := []string{}
	args := []any{}
	argID := 1

	if partialUser.Username != nil {
		setParts = append(setParts, fmt.Sprintf("username=$%d", argID))
		args = append(args, *partialUser.Username)
		argID++
	}
	if partialUser.Password != nil {
		setParts = append(setParts, fmt.Sprintf("password=$%d", argID))
		args = append(args, *partialUser.Password)
		argID++
	}
	if partialUser.HiRes != nil {
		setParts = append(setParts, fmt.Sprintf("hi_res=$%d", argID))
		args = append(args, *partialUser.HiRes)
		argID++
	}
	if partialUser.Role != nil {
		setParts = append(setParts, fmt.Sprintf("role=$%d", argID))
		args = append(args, *partialUser.Role)
		argID++
	}
	args = append(args, partialUser.ID)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d RETURNING *", strings.Join(setParts, ", "), argID)

	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return model.User{}, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return model.User{}, err
	}

	return user.ToUser(), nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	if _, err := tx.Exec(ctx, "DELETE FROM users WHERE id = $1", userID); err != nil {
		return dto.Error(err)
	}
	return nil
}
