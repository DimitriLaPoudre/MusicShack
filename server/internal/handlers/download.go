package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/gin-gonic/gin"
)

func AddDownload(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.RequestDownload
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Type {
	case "song":
		services.DownloadManager.AddSong(userId, req.Provider, req.Id)
	case "album":
		services.DownloadManager.AddAlbum(userId, req.Provider, req.Id)
	case "artist":
		services.DownloadManager.AddArtist(userId, req.Provider, req.Id)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func ListDownload(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks := services.DownloadManager.List(userId)
	c.JSON(http.StatusOK, tasks)
}

func DeleteDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Remove(userId, taskId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RetryDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Retry(userId, taskId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func CancelDownload(c *gin.Context) {
	taskIdRaw := c.Param("id")
	taskIdBadType, err := strconv.ParseUint(taskIdRaw, 10, 0)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	taskId := uint(taskIdBadType)
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = services.DownloadManager.Cancel(userId, taskId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RetryDownloads(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.DownloadManager.RetryAll(userId)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func DoneDownloads(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	services.DownloadManager.Done(userId)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
