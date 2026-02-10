package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
	"github.com/gin-gonic/gin"
)

func UploadSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info := models.MetadataInfo{
		Explicit: "false",
	}
	info.Title = strings.ReplaceAll(c.Request.FormValue("title"), "/", "_")
	if info.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title field empty"})
		return
	}
	info.Album = strings.ReplaceAll(c.Request.FormValue("album"), "/", "_")
	if info.Album == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album field empty"})
		return
	}
	albumArtistsRaw := strings.ReplaceAll(c.Request.FormValue("artists"), "/", "_")
	if albumArtistsRaw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
	}
	info.AlbumArtists = strings.Split(albumArtistsRaw, ",")
	for i, artist := range info.AlbumArtists {
		info.AlbumArtists[i] = strings.TrimSpace(artist)
	}
	artistsRaw := strings.ReplaceAll(c.Request.FormValue("artists"), "/", "_")
	if artistsRaw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
	}
	info.Artists = strings.Split(artistsRaw, ",")
	for i, artist := range info.Artists {
		info.Artists[i] = strings.TrimSpace(artist)
	}
	if value := c.Request.FormValue("trackNumber"); value != "" {
		info.TrackNumber = value
	}
	if value := c.Request.FormValue("volumeNumber"); value != "" {
		info.VolumeNumber = value
	}
	if value := c.Request.FormValue("isrc"); value != "" {
		info.Isrc = value
	} else {
		if random, err := utils.GenerateRandomString(13); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			info.Isrc = random
		}
	}
	if value := c.Request.FormValue("releaseDate"); value != "" {
		info.ReleaseDate = value
	}
	if c.Request.FormValue("explicit") == "on" {
		info.Explicit = "true"
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	filename := header.Filename
	fmt.Println("filename:", filename)
	extensionIndex := strings.LastIndex(filename, ".")
	if extensionIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file don't have an extension"})
		return
	}

	err = services.UploadLibrarySong(c.Request.Context(), userId, file, filename[extensionIndex+1:], info)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func EditSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id uint
	if result, err := strconv.ParseUint(c.Param("id"), 10, 0); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		id = uint(result)
	}

	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var edit models.RequestEditSong
	if err := c.ShouldBindJSON(&edit); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := metadata.ApplyMetadata(filepath.Join(userPath, song.Path), models.MetadataInfo(edit)); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//prevoir le rename

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func GetSongCover(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id uint
	if result, err := strconv.ParseUint(c.Param("id"), 10, 0); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		id = uint(result)
	}

	img, err := services.GetLibrarySongCover(userId, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/jpg", img)
}

func ListSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var limit int
	if result := c.Query("limit"); result == "" {
		limit = 10
	} else {
		if result, err := strconv.Atoi(result); err != nil {
			limit = 10
		} else {
			limit = result
		}
	}

	var offset int
	if result := c.Query("offset"); result == "" {
		offset = 0
	} else {
		if result, err := strconv.Atoi(result); err != nil {
			offset = 0
		} else {
			offset = result
		}
	}

	q := c.Query("q")
	q = strings.ReplaceAll(q, "_", "\\_")
	q = strings.ReplaceAll(q, "/", "\\_")

	total, err := repository.CountSongByUserID(userId, q)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	songs, err := repository.ListSongByUserID(userId, q, limit, offset)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	list := make([]models.ResponseSong, 0)
	for _, dbSong := range songs {
		if song, err := services.GetLibrarySong(dbSong); err != nil {
			log.Println("ListSong:", dbSong, ":", err)
			continue
		} else {
			list = append(list, song)
		}
	}

	c.JSON(http.StatusOK, models.ResponseLibrary{Total: int(total), Count: len(list), Limit: limit, Offset: offset, Items: list})
}

func DeleteSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id uint
	if result, err := strconv.ParseUint(c.Param("id"), 10, 0); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		id = uint(result)
	}

	if err := services.DeleteLibrarySong(userId, id); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SyncLibrary(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.SyncUserLibrary(userId); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
