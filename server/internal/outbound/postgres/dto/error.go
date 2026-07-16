package dto

import (
	"errors"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

func Error(err error) error {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return model.ErrConflict

			case "23503": // foreign_key_violation
				return model.ErrInvalidInput

			case "23502": // not_null_violation
				return model.ErrInvalidInput

			default:
				return model.ErrInternal
			}
		}
		return err
	}
	return nil
}
