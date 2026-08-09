package model

import (
	"time"

	"github.com/google/uuid"
)

type Song struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ISRC      string    `db:"isrc"`
	Path      string    `db:"path"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SongFilter struct {
	ID     *uuid.UUID
	UserID *uuid.UUID
}
