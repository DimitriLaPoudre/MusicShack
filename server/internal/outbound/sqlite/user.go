package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/outbound/sqlite/dto"
	"github.com/google/uuid"
)

func (r *SQLiteRepository) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	tx := r.getTx(ctx)

	query := "INSERT INTO users (id, username, password, hi_res, role) VALUES (?, ?, ?, ?, ?) RETURNING *"
	query = tx.Rebind(query)

	dbUser := dto.User{}
	err := tx.GetContext(ctx,
		&dbUser,
		query,
		user.ID, user.Username, user.Password, user.HiRes, user.Role)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", dto.Error(err))
	}

	return dbUser.ToUser(), nil
}

func (r *SQLiteRepository) GetUserByFilter(ctx context.Context, filter model.UserFilter) (model.User, error) {
	tx := r.getTx(ctx)

	setParts := []string{}
	args := []any{}

	if filter.ID != nil {
		setParts = append(setParts, "id=?")
		args = append(args, *filter.ID)
	}
	if filter.Username != nil {
		setParts = append(setParts, "username=?")
		args = append(args, *filter.Username)
	}
	if filter.Password != nil {
		setParts = append(setParts, "password=?")
		args = append(args, *filter.Password)
	}
	if filter.HiRes != nil {
		setParts = append(setParts, "hi_res=?")
		args = append(args, *filter.HiRes)
	}
	if filter.Role != nil {
		setParts = append(setParts, "role=?")
		args = append(args, *filter.Role)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM users LIMIT 1"
	} else {
		query = fmt.Sprintf("SELECT * FROM users WHERE %s LIMIT 1", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbUser := dto.User{}
	err := tx.GetContext(ctx, &dbUser, query, args...)
	if err != nil {
		return model.User{}, fmt.Errorf("get user with filter %v: %w", filter, dto.Error(err))
	}

	return dbUser.ToUser(), nil
}

func (r *SQLiteRepository) ListUsersByFilter(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	tx := r.getTx(ctx)

	setParts := []string{}
	args := []any{}

	if filter.ID != nil {
		setParts = append(setParts, "id=?")
		args = append(args, *filter.ID)
	}
	if filter.Username != nil {
		setParts = append(setParts, "username=?")
		args = append(args, *filter.Username)
	}
	if filter.Password != nil {
		setParts = append(setParts, "password=?")
		args = append(args, *filter.Password)
	}
	if filter.HiRes != nil {
		setParts = append(setParts, "hi_res=?")
		args = append(args, *filter.HiRes)
	}
	if filter.Role != nil {
		setParts = append(setParts, "role=?")
		args = append(args, *filter.Role)
	}

	var query string
	if len(args) == 0 {
		query = "SELECT * FROM users"
	} else {
		query = fmt.Sprintf("SELECT * FROM users WHERE %s", strings.Join(setParts, " AND "))
	}
	query = tx.Rebind(query)

	dbUsers := []dto.User{}
	err := tx.SelectContext(ctx, &dbUsers, query, args...)
	if err != nil {
		return []model.User{}, fmt.Errorf("list user with filter %v: %w", filter, dto.Error(err))
	}

	return dto.UsersToUsers(dbUsers), nil
}

func (r *SQLiteRepository) UpdateUser(ctx context.Context, partialUser model.PartialUser) (model.User, error) {
	tx := r.getTx(ctx)

	setParts := []string{}
	args := []any{}

	if partialUser.Username != nil {
		setParts = append(setParts, "username=?")
		args = append(args, *partialUser.Username)
	}
	if partialUser.Password != nil {
		setParts = append(setParts, "password=?")
		args = append(args, *partialUser.Password)
	}
	if partialUser.HiRes != nil {
		setParts = append(setParts, "hi_res=?")
		args = append(args, *partialUser.HiRes)
	}
	if partialUser.Role != nil {
		setParts = append(setParts, "role=?")
		args = append(args, *partialUser.Role)
	}
	args = append(args, partialUser.ID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id=? RETURNING *", strings.Join(setParts, ", "))
	query = tx.Rebind(query)

	dbUser := dto.User{}
	err := tx.GetContext(ctx, &dbUser, query, args...)
	if err != nil {
		return model.User{}, fmt.Errorf("update user %s: %w", partialUser.ID, dto.Error(err))
	}

	return dbUser.ToUser(), nil
}

func (r *SQLiteRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx := r.getTx(ctx)

	query := "DELETE FROM users WHERE id = ?"
	query = tx.Rebind(query)

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("delete user %s: %w", userID, dto.Error(err))
	}

	return nil
}
