package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
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
	info.Title = c.Request.FormValue("title")
	if info.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title field empty"})
		return
	}
	info.Album = c.Request.FormValue("album")
	if info.Album == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "album field empty"})
		return
	}
	albumArtistsRaw := c.Request.FormValue("artists")
	info.AlbumArtists = strings.Split(albumArtistsRaw, ",")
	for i, artist := range info.AlbumArtists {
		info.AlbumArtists[i] = strings.TrimSpace(artist)
	}
	if len(info.AlbumArtists) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
	}
	artistsRaw := c.Request.FormValue("artists")
	info.Artists = strings.Split(artistsRaw, ",")
	for i, artist := range info.Artists {
		info.Artists[i] = strings.TrimSpace(artist)
	}
	if len(info.Artists) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
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
			info.Isrc = strings.ToUpper(random)
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

	var editRaw models.RequestEditSong
	if err := c.ShouldBindJSON(&editRaw); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	edit := models.MetadataInfo{
		Title:        editRaw.Title,
		ReleaseDate:  editRaw.ReleaseDate,
		TrackNumber:  strconv.FormatUint(uint64(editRaw.TrackNumber), 10),
		VolumeNumber: strconv.FormatUint(uint64(editRaw.VolumeNumber), 10),
		Explicit:     strconv.FormatBool(editRaw.Explicit),
		Isrc:         editRaw.Isrc,
		Album:        editRaw.Album,
		AlbumArtists: editRaw.AlbumArtists,
		Artists:      editRaw.Artists,
		AlbumGain:    strconv.FormatFloat(editRaw.AlbumGain, 'f', 6, 64),
		AlbumPeak:    strconv.FormatFloat(editRaw.AlbumPeak, 'f', 6, 64),
		TrackGain:    strconv.FormatFloat(editRaw.TrackGain, 'f', 6, 64),
		TrackPeak:    strconv.FormatFloat(editRaw.TrackPeak, 'f', 6, 64),
	}

	if edit.Title == "" {
		err := errors.New("title field empty")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(edit.AlbumArtists) == 0 {
		err := errors.New("albumArtists field empty")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(edit.Artists) == 0 {
		err := errors.New("artists field empty")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if edit.Isrc == "" {
		err := errors.New("isrc field empty")
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

	extensionIndex := strings.LastIndex(song.Path, ".")
	if extensionIndex == -1 {
		err := errors.New("file don't have an extension")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newSongPath := filepath.Join(strings.ReplaceAll(edit.AlbumArtists[0], "/", "_"), strings.ReplaceAll(edit.Album, "/", "_"),
		fmt.Sprintf("%s - %s.%s", strings.ReplaceAll(edit.TrackNumber, "/", "_"), strings.ReplaceAll(edit.Title, "/", "_"), song.Path[extensionIndex+1:]))
	if song.Path != newSongPath {
		root, err := os.OpenRoot(userPath)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		defer root.Close()

		if err := root.Rename(song.Path, newSongPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		song.Path = newSongPath
	}

	song.Isrc = edit.Isrc

	// update the song even without change for updatedAt field
	if err := repository.UpdateSongByUserID(userId, song); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := metadata.ApplyMetadata(filepath.Join(userPath, song.Path), models.MetadataInfo(edit)); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
