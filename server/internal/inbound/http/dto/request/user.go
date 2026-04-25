package request

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
)

type CreateUser struct {
	Username string `json:"name" binding:"required,min=3,max=20,alphanumunicode|contains=_"`
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

// type UpdateUser struct {
// 	Name     *string         `json:"name" binding:"required,min=3,max=20,alphanumunicode|contains=_"`
// 	Email    *string         `json:"email" binding:"omitempty,email"`
// 	Password *string         `json:"password"`
// 	Role     *model.UserRole `json:"role"`
// }
//
// func (req *UpdateUser) IntoPartialUser(idStr string) (model.PartialUser, error) {
// 	id, err := IntoUUID(idStr)
// 	if err != nil {
// 		return model.PartialUser{}, err
// 	}
//
// 	var bytePassword []byte
// 	if req.Password != nil {
// 		bytePassword = []byte(*(req.Password))
// 	}
//
// 	return model.PartialUser{
// 		ID:       id,
// 		Name:     req.Name,
// 		Email:    req.Email,
// 		Password: &bytePassword,
// 		Role:     req.Role,
// 	}, nil
// }
