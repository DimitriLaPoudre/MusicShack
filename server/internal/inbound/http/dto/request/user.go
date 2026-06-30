package request

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
)

type CreateUser struct {
	Username string `json:"username" binding:"required,min=3,max=20,alphanumunicode|contains=_"`
	Password string `json:"password" binding:"required"`
	HiRes    bool   `json:"hi_res"`
}

func (req *CreateUser) IntoUser() (model.User, error) {
	return model.User{
		Username: req.Username,
		Password: req.Password,
		HiRes:    req.HiRes,
	}, nil
}

func ParseUserID(str string) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

type UpdateUser struct {
	Username *string `json:"username" binding:"min=3,max=20,alphanumunicode|contains=_"`
	Password *string `json:"password"`
	HiRes    *bool   `json:"hi_res"`
}

func (req *UpdateUser) IntoPartialUser(str string) (model.PartialUser, error) {
	id, err := ParseUserID(str)
	if err != nil {
		return model.PartialUser{}, err
	}

	return model.PartialUser{
		ID:       id,
		Username: req.Username,
		Password: req.Password,
		HiRes:    req.HiRes,
	}, nil
}
