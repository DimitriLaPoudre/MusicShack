package utils

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	slog.ErrorContext(c.Request.Context(), "handler forward domain error", slog.String("err", err.Error()))

	switch {
	case errors.Is(err, model.ErrDownloadInvalidID):
		c.JSON(http.StatusBadRequest, response.NewError(err))
	case errors.Is(err, model.ErrDownloadNotFound),
		errors.Is(err, model.ErrNotFound),
		errors.Is(err, model.ErrUserNotFound),
		errors.Is(err, model.ErrPluginNotFound),
		errors.Is(err, model.ErrPluginDataNotFound):
		c.JSON(http.StatusNotFound, response.NewError(err))
	// case errors.Is(err, model.ErrRoleInvalid):
	// 	c.JSON(http.StatusUnprocessableEntity, response.Error{Message: err.Error()})
	// case errors.Is(err, model.ErrEmailDuplicate):
	// 	c.JSON(http.StatusConflict, response.Error{Message: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, response.NewError(err))
	}
}
