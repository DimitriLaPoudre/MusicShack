package model

import (
	"time"

	"github.com/google/uuid"
)

type Instance struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Provider string
	Plugin   string
	Url      string
	Ping     time.Duration
}
