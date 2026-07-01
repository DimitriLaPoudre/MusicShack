package request

import (
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
)

type CreateInstance struct {
	Url string `json:"url" binding:"required,url"`
}

func (req *CreateInstance) IntoInstance() (model.Instance, error) {
	return model.Instance{
		Url: req.Url,
	}, nil
}
