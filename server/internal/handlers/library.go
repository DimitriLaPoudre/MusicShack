package handlers

import (
	"errors"
	"fmt"
	"log"
	"mime/multipart"
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
		AlbumArtists: make([]string, 0),
		Artists:      make([]string, 0),
		Explicit:     "false",
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
	albumArtists := strings.Split(c.Request.FormValue("albumArtists"), ",")
	for _, artist := range albumArtists {
		artist = strings.TrimSpace(artist)
		if artist != "" {
			info.AlbumArtists = append(info.AlbumArtists, artist)
		}
	}
	if len(info.AlbumArtists) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
	}

	artists := strings.Split(c.Request.FormValue("artists"), ",")
	for _, artist := range artists {
		artist = strings.TrimSpace(artist)
		if artist != "" {
			info.Artists = append(info.Artists, artist)
		}
	}
	if len(info.Artists) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "artists field empty"})
		return
	}

	info.TrackNumber = c.Request.FormValue("trackNumber")
	info.VolumeNumber = c.Request.FormValue("volumeNumber")
	info.ReleaseDate = c.Request.FormValue("releaseDate")

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

	if c.Request.FormValue("explicit") == "on" {
		info.Explicit = "true"
	}

	cover, _, err := c.Request.FormFile("cover")
	if err != nil {
		if err != http.ErrMissingFile {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		defer cover.Close()
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	filename := fileHeader.Filename
	extensionIndex := strings.LastIndex(filename, ".")
	if extensionIndex == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file don't have an extension"})
		return
	}

	err = services.UploadLibrarySong(c.Request.Context(), userId, file, filename[extensionIndex+1:], info, cover)
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

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refTags, err := metadata.ReadTags(filepath.Join(userPath, song.Path))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var edit models.RequestEditSong
	if err := c.ShouldBind(&edit); err != nil {
		err := fmt.Errorf("c.ShouldBind: %w", err)
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cover multipart.File
	if edit.Cover != nil {
		if cover, err = edit.Cover.Open(); err != nil {
			err := fmt.Errorf("edit.Cover.Open: %w", err)
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var title string
	var album string
	var artist string
	var trackNumber string
	tags := map[string][]string{}
	if edit.Title != nil {
		if *edit.Title == "" {
			err := errors.New("title field empty")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tags[metadata.Title] = []string{*edit.Title}
			title = *edit.Title
		}
	} else {
		title = refTags[metadata.Title][0]
	}
	if edit.Album != nil {
		if *edit.Album == "" {
			err := errors.New("album field empty")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tags[metadata.Album] = []string{*edit.Album}
			album = *edit.Album
		}
	} else {
		album = refTags[metadata.Album][0]
	}
	if edit.AlbumArtists != nil {
		if len(*edit.AlbumArtists) == 0 {
			err := errors.New("albumArtists field empty")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tags[metadata.AlbumArtists] = *edit.AlbumArtists
			artist = (*edit.AlbumArtists)[0]
		}
	} else {
		artist = refTags[metadata.AlbumArtists][0]
	}
	if edit.Artists != nil {
		if len(*edit.Artists) == 0 {
			err := errors.New("artists field empty")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tags[metadata.Artists] = *edit.Artists
		}
	}
	if edit.ReleaseDate != nil {
		tags[metadata.ReleaseDate] = []string{*edit.ReleaseDate}
	}
	if edit.TrackNumber != nil {
		tags[metadata.TrackNumber] = []string{strconv.FormatUint(uint64(*edit.TrackNumber), 10)}
		trackNumber = strconv.FormatUint(uint64(*edit.TrackNumber), 10)
	} else {
		trackNumber = refTags[metadata.TrackNumber][0]
	}
	if edit.VolumeNumber != nil {
		tags[metadata.VolumeNumber] = []string{strconv.FormatUint(uint64(*edit.VolumeNumber), 10)}
	}
	if edit.Explicit != nil {
		tags[metadata.Explicit] = []string{strconv.FormatBool(*edit.Explicit)}
	}
	if edit.AlbumGain != nil {
		tags[metadata.AlbumGain] = []string{strconv.FormatFloat(*edit.AlbumGain, 'f', 6, 64)}
	}
	if edit.AlbumPeak != nil {
		tags[metadata.AlbumPeak] = []string{strconv.FormatFloat(*edit.AlbumPeak, 'f', 6, 64)}
	}
	if edit.TrackGain != nil {
		tags[metadata.TrackGain] = []string{strconv.FormatFloat(*edit.TrackGain, 'f', 6, 64)}
	}
	if edit.TrackPeak != nil {
		tags[metadata.TrackPeak] = []string{strconv.FormatFloat(*edit.TrackPeak, 'f', 6, 64)}
	}
	if edit.Isrc != nil {
		if *edit.Isrc == "" {
			err := errors.New("isrc field empty")
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {
			tags[metadata.ISRC] = []string{*edit.Isrc}
			song.Isrc = *edit.Isrc
		}
	}

	extensionIndex := strings.LastIndex(song.Path, ".")
	if extensionIndex == -1 {
		err := errors.New("file don't have an extension")
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	newSongPath := filepath.Join(strings.ReplaceAll(artist, "/", "_"), strings.ReplaceAll(album, "/", "_"),
		fmt.Sprintf("%s - %s.%s", strings.ReplaceAll(trackNumber, "/", "_"), strings.ReplaceAll(title, "/", "_"), song.Path[extensionIndex+1:]))
	if song.Path != newSongPath {
		root, err := os.OpenRoot(userPath)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		defer root.Close()

		oldSongPath := song.Path
		if err := utils.RootDuplicateFile(root, oldSongPath, newSongPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := metadata.WriteTags(filepath.Join(userPath, newSongPath), tags, false); err != nil {
			_ = root.Remove(newSongPath)
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			_ = root.Remove(oldSongPath)
		}
		if cover != nil {
			if err := metadata.WriteCover(filepath.Join(userPath, newSongPath), cover); err != nil {
				_ = root.Remove(newSongPath)
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				_ = root.Remove(oldSongPath)
			}
		}

		song.Path = newSongPath
	} else {
		if err := metadata.WriteTags(filepath.Join(userPath, newSongPath), tags, false); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if cover != nil {
			if err := metadata.WriteCover(filepath.Join(userPath, newSongPath), cover); err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// update the song even without change for updatedAt field
	if err := repository.UpdateSongByUserID(userId, song); err != nil {
		log.Println(err)
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
