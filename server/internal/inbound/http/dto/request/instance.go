package request

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
)

type CreateInstance struct {
	Url string `json:"url" binding:"required,url"`
}

func (req *CreateInstance) IntoInstance(userID uuid.UUID) (model.Instance, error) {
	return model.Instance{
		UserID: userID,
		Url:    req.Url,
	}, nil
}
