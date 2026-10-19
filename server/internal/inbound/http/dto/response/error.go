package response

import (
	"errors"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type Error struct {
	Message string `json:"error"`
}

func NewError(err error) Error {
	if err, ok := errors.AsType[*model.Error](err); ok {
		return Error{
			Message: err.Message(),
		}
	}

	return Error{
		Message: err.Error(),
	}
}
