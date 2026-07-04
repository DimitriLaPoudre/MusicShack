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
	me, err := utils.GetFromContext[model.User](c, "me")
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

func (h *PluginHandler) GetAlbum(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	album, err := h.plugin.GetAlbum(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.AlbumToResponse(album)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetArtist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	artist, err := h.plugin.GetArtist(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.ArtistToResponse(artist)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetPlaylist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	playlist, err := h.plugin.GetPlaylist(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err, h.l)
		return
	}

	resp := response.PlaylistToResponse(playlist)
	c.JSON(http.StatusOK, resp)
}
