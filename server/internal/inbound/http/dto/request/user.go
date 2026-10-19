package request

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"uuid"
)

type CreateUser struct {
	Username string `json:"username" binding:"required,min=3,max=32,regex=^[a-zA-Z0-9_]+$"`
	Password string `json:"password" binding:"required,min=8,max=128,regex=^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])[a-zA-Z0-9]+$"`
	HiRes    bool   `json:"hi_res"`
	Role     string `json:"role"`
}

func (req CreateUser) IntoUser() (model.User, error) {
	return model.User{
		Username: req.Username,
		Password: req.Password,
		HiRes:    req.HiRes,
		Role:     model.UserRole(req.Role),
	}, nil
}

type UpdateUser struct {
	Username *string `json:"username" binding:"omitempty,min=3,max=32,regex=^[a-zA-Z0-9_]+$"`
	Password *string `json:"password" binding:"omitempty,min=8,max=128,regex=^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])[a-zA-Z0-9]+$"`
	HiRes    *bool   `json:"hi_res"`
}

func (req UpdateUser) IntoPartialUser(id uuid.UUID) model.PartialUser {
	return model.PartialUser{
		ID:       id,
		Username: req.Username,
		Password: req.Password,
		HiRes:    req.HiRes,
	}
}
