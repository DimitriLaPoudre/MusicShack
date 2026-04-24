package dto

import (
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
)

type Admin struct {
	ID        int       `db:"id"`
	Password  string    `db:"password"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (u *Admin) ToAdmin() *model.Admin {
	return &model.Admin{
		Password:  u.Password,
		Token:     u.Token,
		ExpiresAt: u.ExpiresAt,
	}
}
