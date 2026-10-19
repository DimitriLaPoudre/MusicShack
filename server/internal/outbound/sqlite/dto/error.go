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
			return model.NewOutboundError(model.ErrNotFound, err)
		}

		if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok {
			switch sqliteErr.Code() {
			case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
				return model.NewOutboundError(model.NewError(model.ErrConflict, model.MessageCodeAlreadyExists, "a field"), err)

			case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
				return model.NewOutboundError(model.NewError(model.ErrInvalid, model.MessageCodeInvalid, "a field"), err)

			case sqlite3.SQLITE_CONSTRAINT_NOTNULL:
				return model.NewOutboundError(model.NewError(model.ErrInvalid, model.MessageCodeInvalid, "a field"), err)

			default:
				return model.NewOutboundError(model.NewError(model.ErrInternal, model.MessageCodeInternalError), err)
			}
		}
	}

	return err
}
