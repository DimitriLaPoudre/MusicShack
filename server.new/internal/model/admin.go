package model

import "time"

type Admin struct {
	Password  string
	Token     string
	ExpiresAt time.Time
}
