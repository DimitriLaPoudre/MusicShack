package handler

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	download *service.DownloadService
}

func NewDownloadHandler(download *service.DownloadService) DownloadHandler {
	return DownloadHandler{
		download: download,
	}
}

func (h *DownloadHandler) Download(c *gin.Context) {
}

func (h *DownloadHandler) Retry(c *gin.Context) {
}

func (h *DownloadHandler) RetryAll(c *gin.Context) {
}

func (h *DownloadHandler) Cancel(c *gin.Context) {
}

func (h *DownloadHandler) RemoveDone(c *gin.Context) {
}

func (h *DownloadHandler) Remove(c *gin.Context) {
}

func (h *DownloadHandler) List(c *gin.Context) {
}
