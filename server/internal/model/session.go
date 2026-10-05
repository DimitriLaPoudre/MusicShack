package model

import (
	"time"

	"uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type SessionFilter struct {
	ID     *uuid.UUID
	UserID *uuid.UUID
	Token  *string
}
