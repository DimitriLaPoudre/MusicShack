package response

import "github.com/Ascension-EIP/Ascension/apps/server/internal/model"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	HiRes    bool   `json:"hi_res"`
}

func UserToResponse(user *model.User) *User {
	return &User{
		ID:       user.ID.String(),
		Username: user.Username,
		HiRes:    user.HiRes,
	}
}

func UsersToResponse(users []*model.User) []*User {
	r := []*User{}
	for _, user := range users {
		r = append(r, UserToResponse(user))
	}
	return r
}
