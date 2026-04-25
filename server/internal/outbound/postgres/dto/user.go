package dto

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	HiRes    bool      `db:"hi_res"`
}

func (u *User) ToUser() *model.User {
	return &model.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		HiRes:    u.HiRes,
	}
}

func UsersToUsers(dto []*User) []*model.User {
	users := []*model.User{}
	for _, user := range dto {
		users = append(users, &model.User{
			ID:       user.ID,
			Username: user.Username,
			Password: user.Password,
			HiRes:    user.HiRes,
		})
	}
	return users
}
