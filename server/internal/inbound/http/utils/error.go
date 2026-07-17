package utils

import (
	"log/slog"
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	slog.ErrorContext(c.Request.Context(), "handler forward domain error", slog.String("err", err.Error()))

	switch {
	// case errors.Is(err, model.ErrRoleInvalid):
	// 	c.JSON(http.StatusUnprocessableEntity, response.Error{Message: err.Error()})
	// case errors.Is(err, model.ErrEmailDuplicate):
	// 	c.JSON(http.StatusConflict, response.Error{Message: err.Error()})
	// case errors.Is(err, model.ErrUserNotFound):
	// 	c.JSON(http.StatusNotFound, response.Error{Message: err.Error()})
	default:
		// _ = c.Error(err)
		c.JSON(http.StatusInternalServerError, response.NewError(err))
	}
}
