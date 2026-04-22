package model

import "time"

type Admin struct {
	ID        bool
	Password  string
	Token     string
	ExpiresAt time.Time
}
