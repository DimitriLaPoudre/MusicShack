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

func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	if user == nil {
		return nil, model.ErrUnknown
	}
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"INSERT INTO users (id, username, password, hi_res, role) VALUES ($1, $2, $3, $4, $5) RETURNING *",
		user.ID, user.Username, user.Password, user.HiRes, user.Role)
	if err != nil {
		return nil, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM users WHERE id = $1 LIMIT 1",
		userID)
	if err != nil {
		return nil, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM users WHERE username = $1 LIMIT 1",
		username)
	if err != nil {
		return nil, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) GetUserWithFilter(ctx context.Context, filter *model.FilterUser) (*model.User, error) {
	if filter == nil {
		return nil, model.ErrUnknown
	}

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
		return nil, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) ListUsersWithFilter(ctx context.Context, filter *model.FilterUser) ([]*model.User, error) {
	if filter == nil {
		return nil, model.ErrUnknown
	}

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
		return nil, err
	}

	dbUsers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}

	return dto.UsersToUsers(dbUsers), nil
}

func (r *PostgresRepository) ListAllUsers(ctx context.Context) ([]*model.User, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM users ORDER BY username ASC")
	if err != nil {
		return nil, err
	}

	dbUsers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
	}
	return dto.UsersToUsers(dbUsers), nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, partialUser *model.PartialUser) (*model.User, error) {
	if partialUser == nil {
		return nil, model.ErrUnknown
	}

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
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err := dto.Error(err); err != nil {
		return nil, err
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
