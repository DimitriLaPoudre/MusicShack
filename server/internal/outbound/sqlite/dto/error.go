package dto

import (
	"database/sql"
	"errors"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func Error(err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNotFound
		}

		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return model.ErrConflict

			case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
				return model.ErrInvalidInput

			case sqlite3.SQLITE_CONSTRAINT_NOTNULL:
				return model.ErrInvalidInput

			default:
				return model.ErrInternal
			}
		}
	}

	return err
}
