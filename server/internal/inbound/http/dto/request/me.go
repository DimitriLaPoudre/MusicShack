package request

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type UpdateMe struct {
	Username *string `json:"username" binding:"min=3,max=20,alphanumunicode|contains=_"`
	Password *string `json:"password"`
	HiRes    *bool   `json:"hi_res"`
}

func (req UpdateMe) IntoPartialUser(id uuid.UUID) (model.PartialUser, error) {
	return model.PartialUser{
		ID:       id,
		Username: req.Username,
		Password: req.Password,
		HiRes:    req.HiRes,
	}, nil
}
