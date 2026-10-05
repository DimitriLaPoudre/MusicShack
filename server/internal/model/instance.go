package model

import (
	"time"

	"uuid"
)

type Instance struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Provider string
	Plugin   string
	URL      string
	Ping     *time.Duration
}

type InstanceFilter struct {
	ID       *uuid.UUID
	UserID   *uuid.UUID
	Provider *string
	Plugin   *string
	URL      *string
}
