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
	Role     string    `db:"role"`
}

func (u *User) ToUser() *model.User {
	return &model.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		HiRes:    u.HiRes,
		Role:     model.UserRole(u.Role),
	}
}

func UsersToUsers(dto []*User) []*model.User {
	users := []*model.User{}
	for _, u := range dto {
		users = append(users, &model.User{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
			HiRes:    u.HiRes,
			Role:     model.UserRole(u.Role),
		})
	}
	return users
}
