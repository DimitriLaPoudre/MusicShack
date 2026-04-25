package postgres

import (
	"context"

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
		"INSERT INTO users (id, username, password, hi_res) VALUES ($1, $2, $3, $4) RETURNING *",
		user.ID, user.Username, user.Password, user.HiRes)
	if err != nil {
		return nil, err
	}

	dbUser, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		return nil, err
	}

	return dbUser.ToUser(), nil
}

func (r *PostgresRepository) ListAllUsers(ctx context.Context) ([]*model.User, error) {
	tx := r.getTx(ctx)

	rows, err := tx.Query(ctx,
		"SELECT * FROM users ORDER BY username ASC")
	if err != nil {
		return nil, err
	}

	dbUsers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[dto.User])
	if err != nil {
		return nil, err
	}
	return dto.UsersToUsers(dbUsers), nil
}

// func (r *PostgresRepository) UpdateUser(ctx context.Context, partialUser *model.PartialUser) (*model.User, error) {
// 	if partialUser == nil {
// 		return nil, model.ErrUnknown
// 	}
//
// 	setParts := []string{}
// 	args := []any{}
// 	argID := 1
//
// 	if partialUser.Name != nil {
// 		setParts = append(setParts, fmt.Sprintf("name=$%d", argID))
// 		args = append(args, *partialUser.Name)
// 		argID++
// 	}
// 	if partialUser.Email != nil {
// 		setParts = append(setParts, fmt.Sprintf("email=$%d", argID))
// 		args = append(args, *partialUser.Email)
// 		argID++
// 	}
// 	if partialUser.Password != nil {
// 		setParts = append(setParts, fmt.Sprintf("password=$%d", argID))
// 		args = append(args, *partialUser.Password)
// 		argID++
// 	}
// 	if partialUser.Role != nil {
// 		setParts = append(setParts, fmt.Sprintf("role=$%d", argID))
// 		args = append(args, *partialUser.Role)
// 		argID++
// 	}
// 	args = append(args, partialUser.ID)
// 	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d RETURNING *", strings.Join(setParts, ", "), argID)
//
// 	tx := r.getTx(ctx)
//
// 	rows, err := tx.Query(context.Background(), query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[dto.User])
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return user.ToUser(), nil
// }

func (r *PostgresRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	_, err := tx.Exec(ctx,
		"DELETE FROM users WHERE id = $1 LIMIT 1",
		userID)
	if err != nil {
		return err
	}
	return nil
}
