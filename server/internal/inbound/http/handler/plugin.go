package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/utils"
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/gin-gonic/gin"
)

type PluginHandler struct {
	plugin *service.PluginService
}

func NewPluginHandler(plugin *service.PluginService) PluginHandler {
	return PluginHandler{
		plugin: plugin,
	}
}

func (h *PluginHandler) GetSong(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	song, err := h.plugin.GetSong(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.SongToResponse(song)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetAlbum(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	album, err := h.plugin.GetAlbum(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.AlbumToResponse(album)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetArtist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	artist, err := h.plugin.GetArtist(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.ArtistToResponse(artist)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetPlaylist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param("provider")
	id := c.Param("id")

	playlist, err := h.plugin.GetPlaylist(c.Request.Context(), me, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.PlaylistToResponse(playlist)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetSearch(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, "me")
	if err != nil {
		utils.Error(c, err)
		return
	}

	q := c.Query("q")

	results, err := h.plugin.Search(c.Request.Context(), me, q)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.SearchResultToResponse(results)
	c.JSON(http.StatusOK, resp)
}
