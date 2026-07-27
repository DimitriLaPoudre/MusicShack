package handler

import (
	"net/http"

	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/request"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/dto/response"
	"github.com/DimitriLaPoudre/MusicShack/internal/inbound/http/macro"
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

func (h *PluginHandler) GetSongInfo(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	songInfo, err := h.plugin.GetSongInfo(c.Request.Context(), me.ID, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.SongInfoToResponse(songInfo)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetAlbumInfo(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	albumInfo, err := h.plugin.GetAlbumInfo(c.Request.Context(), me.ID, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.AlbumInfoToResponse(albumInfo)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetAlbumSongs(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	var query request.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	albumSongs, err := h.plugin.GetAlbumSongs(c.Request.Context(), me.ID, provider, id, query.Limit, query.Offset)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.AlbumSongsToResponse(albumSongs)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetArtistInfo(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	artistInfo, err := h.plugin.GetArtistInfo(c.Request.Context(), me.ID, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.ArtistInfoToResponse(artistInfo)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetArtistAlbums(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	var query request.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	artistAlbums, err := h.plugin.GetArtistAlbums(c.Request.Context(), me.ID, provider, id, query.Limit, query.Offset)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.ArtistAlbumsToResponse(artistAlbums)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetPlaylistInfo(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	playlistInfo, err := h.plugin.GetPlaylistInfo(c.Request.Context(), me.ID, provider, id)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.PlaylistInfoToResponse(playlistInfo)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetPlaylistSongs(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	provider := c.Param(macro.Provider)
	id := c.Param(macro.ID)

	var query request.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	playlist, err := h.plugin.GetPlaylistSongs(c.Request.Context(), me.ID, provider, id, query.Limit, query.Offset)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.PlaylistSongsToResponse(playlist)
	c.JSON(http.StatusOK, resp)
}

func (h *PluginHandler) GetSearch(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var query request.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	results, err := h.plugin.Search(c.Request.Context(), me.ID, query.Q, query.Limit, query.Offset)
	if err != nil {
		utils.Error(c, err)
		return
	}

	resp := response.SearchResultToResponse(results)
	c.JSON(http.StatusOK, resp)
}
