package request

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type CreateInstance struct {
	URL string `json:"url" binding:"required,url"`
}

func (req CreateInstance) IntoInstance(userID uuid.UUID) (model.Instance, error) {
	return model.Instance{
		UserID: userID,
		URL:    req.URL,
	}, nil
}
