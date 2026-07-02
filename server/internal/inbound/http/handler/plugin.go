package handler

import (
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/dto/response"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/inbound/http/utils"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type PluginHandler struct {
	l      *zerolog.Logger
	plugin *usecase.PluginUseCase
}

func NewPluginHandler(l *zerolog.Logger, plugin *usecase.PluginUseCase) PluginHandler {
	return PluginHandler{
		l:      l,
		plugin: plugin,
	}
}

func (h *PluginHandler) GetSong(c *gin.Context) {
	me, err := utils.GetFromContext[*model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	song, err := h.plugin.GetSong(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.SongToResponse(song)
	c.JSON(http.StatusOK, resp)
}
