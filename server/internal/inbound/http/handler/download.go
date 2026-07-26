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
	"github.com/google/uuid"
)

type DownloadHandler struct {
	download *service.DownloadService
}

func NewDownloadHandler(download *service.DownloadService) DownloadHandler {
	return DownloadHandler{
		download: download,
	}
}

func (h *DownloadHandler) DownloadArtist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateDownload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.download.DownloadArtist(c.Request.Context(), me.ID, req.Provider, req.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) DownloadAlbum(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateDownload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.download.DownloadAlbum(c.Request.Context(), me.ID, req.Provider, req.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) DownloadPlaylist(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateDownload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.download.DownloadPlaylist(c.Request.Context(), me.ID, req.Provider, req.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) DownloadSong(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	var req request.CreateDownload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewError(err))
		return
	}

	if err := h.download.DownloadSong(c.Request.Context(), me.ID, req.Provider, req.ID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) Retry(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	rawTaskID := c.Param(macro.ID)
	taskID, err := uuid.Parse(rawTaskID)
	if err != nil {
		utils.Error(c, model.ErrDownloadInvalidID)
		return
	}

	if err := h.download.Retry(me.ID, taskID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) RetryAll(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.download.RetryAll(me.ID)

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) Cancel(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	rawTaskID := c.Param(macro.ID)
	taskID, err := uuid.Parse(rawTaskID)
	if err != nil {
		utils.Error(c, model.ErrDownloadInvalidID)
		return
	}

	if err := h.download.Cancel(me.ID, taskID); err != nil {
		utils.Error(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) RemoveDone(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	h.download.RemoveDone(me.ID)
	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) Remove(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	rawTaskID := c.Param(macro.ID)
	taskID, err := uuid.Parse(rawTaskID)
	if err != nil {
		utils.Error(c, model.ErrDownloadInvalidID)
		return
	}

	h.download.Remove(me.ID, taskID)
	c.Status(http.StatusNoContent)
}

func (h *DownloadHandler) List(c *gin.Context) {
	me, err := utils.GetFromContext[model.User](c, macro.Me)
	if err != nil {
		utils.Error(c, err)
		return
	}

	tasks := h.download.List(me.ID)

	resp := response.DownloadTasksToResponse(tasks)
	c.JSON(http.StatusOK, resp)
}
