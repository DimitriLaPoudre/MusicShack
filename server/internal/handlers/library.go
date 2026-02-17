package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	var upload models.RequestUploadSong
	if err := c.ShouldBind(&upload); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			fmt.Errorf("c.ShouldBind: %w", err))
		return
	}

	if upload.File == nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			errors.New("RequestUploadSong.file is empty"))
		return
	}

	file, err := upload.File.Open()
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			errors.New("RequestUploadSong.file is empty"))
		return
	}
	defer file.Close()
	extension := filepath.Ext(upload.File.Filename)

	tmpFile, err := utils.CopyTemporary(file)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			fmt.Println(fmt.Errorf("tmpFile.Close: %w", err))
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Println(fmt.Errorf("tmpFile.Close: %w", err))
		}
	}()

	if upload.Cover != nil {
		cover, err := upload.Cover.Open()
		if err != nil {
			utils.GinPrettyError(c, http.StatusBadRequest,
				fmt.Errorf("upload.Cover.Open: %w", err))
			return
		}
		defer cover.Close()

		if err := metadata.WriteCover(tmpFile.Name(), cover); err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}
	}

	newTags := upload.ToTags()
	if err := metadata.WriteTags(tmpFile.Name(), newTags, false); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	tags, err := metadata.ReadTags(tmpFile.Name())
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	path, err := services.GetSongPathByTags(tags, extension)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	var isrc string
	if value, ok := tags[models.TagISRC]; !ok || len(value) <= 0 ||
		!regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}[0-9]{2}[0-9]{5}$`).MatchString(value[0]) {
		utils.GinPrettyError(c, http.StatusBadRequest, errors.New("isrc field empty and file don't provide default"))
		return
	} else {
		isrc = value[0]
	}

	if err := repository.AddSong(models.Song{UserId: userId, Path: path, Isrc: isrc}); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	rootUser, err := os.OpenRoot(userPath)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	defer rootUser.Close()

	if err := rootUser.MkdirAll(filepath.Dir(path), 0755); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	newFile, err := rootUser.Create(path)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	defer newFile.Close()

	if _, err := tmpFile.Seek(0, 0); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			fmt.Errorf("tmpFile.Seek: %w", err))
		return
	}
	if _, err := io.Copy(newFile, tmpFile); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			fmt.Errorf("io.Copy: %w", err))
		return
	}
	if err := newFile.Sync(); err != nil {
		fmt.Println(fmt.Errorf("newFile.Sync: %w", err))
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func EditSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	result, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("strconv.ParseUint: %w", err))
		return
	}
	id := uint(result)

	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	originalFile, err := os.Open(filepath.Join(userPath, song.Path))
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError,
			fmt.Errorf("os.Open: %w", err))
		return
	}
	copyFile, err := utils.CopyTemporary(originalFile)
	originalFile.Close()
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}
	defer func() {
		if copyFile != nil {
			copyFile.Close()
			os.Remove(copyFile.Name())
		}
	}()

	var edit models.RequestEditSong
	if err := c.ShouldBind(&edit); err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("c.ShouldBind: %w", err))
		return
	}

	newTags := edit.ToTags()
	if err := metadata.WriteTags(copyFile.Name(), newTags, false); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	var cover multipart.File
	if edit.Cover != nil {
		if cover, err = edit.Cover.Open(); err != nil {
			utils.GinPrettyError(c, http.StatusBadRequest,
				fmt.Errorf("edit.Cover.Open: %w", err))
			return
		}
		err := metadata.WriteCover(copyFile.Name(), cover)
		cover.Close()
		if err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}
	}

	tags, err := metadata.ReadTags(copyFile.Name())
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	extension := filepath.Ext(originalFile.Name())
	path, err := services.GetSongPathByTags(tags, extension)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	var isrc string
	if value, ok := tags[models.TagISRC]; !ok || len(value) <= 0 ||
		!regexp.MustCompile(`^[A-Z]{2}[A-Z0-9]{3}[0-9]{2}[0-9]{5}$`).MatchString(value[0]) {
		utils.GinPrettyError(c, http.StatusBadRequest, errors.New("isrc field empty and file don't provide default"))
		return
	} else {
		isrc = value[0]
	}

	if err := repository.UpdateSongByUserID(userId, models.Song{ID: song.ID, Path: path, Isrc: isrc}); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	if originalFile.Name() == path {
		if err := utils.RenameForce(copyFile.Name(), path); err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}
	} else {
		if err := utils.RenameSoft(copyFile.Name(), path); err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}
		if err := os.Remove(originalFile.Name()); err != nil {
			fmt.Println(fmt.Errorf("os.Remove: %w", err))
		}
	}
	copyFile = nil
}

func OldEditSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	result, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("strconv.ParseUint: %w", err))
		return
	}
	id := uint(result)

	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	refTags, err := metadata.ReadTags(filepath.Join(userPath, song.Path))
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	var edit models.RequestEditSong
	if err := c.ShouldBind(&edit); err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("c.ShouldBind: %w", err))
		return
	}

	var cover multipart.File
	if edit.Cover != nil {
		if cover, err = edit.Cover.Open(); err != nil {
			utils.GinPrettyError(c, http.StatusBadRequest,
				fmt.Errorf("edit.Cover.Open: %w", err))
			return
		}
		defer cover.Close()
	}

	var title string
	var album string
	var artist string
	var trackNumber string
	tags := map[string][]string{}
	if edit.Title != nil {
		if *edit.Title == "" {
			utils.GinPrettyError(c, http.StatusBadRequest,
				errors.New("title field empty"))
			return
		} else {
			tags[models.TagTitle] = []string{*edit.Title}
			title = *edit.Title
		}
	} else {
		title = refTags[models.TagTitle][0]
	}
	if edit.Album != nil {
		if *edit.Album == "" {
			utils.GinPrettyError(c, http.StatusBadRequest,
				errors.New("album field empty"))
			return
		} else {
			tags[models.TagAlbum] = []string{*edit.Album}
			album = *edit.Album
		}
	} else {
		album = refTags[models.TagAlbum][0]
	}
	if edit.AlbumArtists != nil {
		if len(*edit.AlbumArtists) == 0 {
			utils.GinPrettyError(c, http.StatusBadRequest,
				errors.New("albumArtists field empty"))
			return
		} else {
			tags[models.TagAlbumArtists] = *edit.AlbumArtists
			artist = (*edit.AlbumArtists)[0]
		}
	} else {
		artist = refTags[models.TagAlbumArtists][0]
	}
	if edit.Artists != nil {
		if len(*edit.Artists) == 0 {
			utils.GinPrettyError(c, http.StatusBadRequest,
				errors.New("artists field empty"))
			return
		} else {
			tags[models.TagArtists] = *edit.Artists
		}
	}
	if edit.ReleaseDate != nil {
		tags[models.TagReleaseDate] = []string{*edit.ReleaseDate}
	}
	if edit.TrackNumber != nil {
		tags[models.TagTrackNumber] = []string{strconv.FormatUint(uint64(*edit.TrackNumber), 10)}
		trackNumber = strconv.FormatUint(uint64(*edit.TrackNumber), 10)
	} else {
		trackNumber = refTags[models.TagTrackNumber][0]
	}
	if edit.VolumeNumber != nil {
		tags[models.TagVolumeNumber] = []string{strconv.FormatUint(uint64(*edit.VolumeNumber), 10)}
	}
	if edit.Explicit != nil {
		tags[models.TagExplicit] = []string{strconv.FormatBool(*edit.Explicit)}
	}
	if edit.AlbumGain != nil {
		tags[models.TagAlbumGain] = []string{strconv.FormatFloat(*edit.AlbumGain, 'f', 6, 64)}
	}
	if edit.AlbumPeak != nil {
		tags[models.TagAlbumPeak] = []string{strconv.FormatFloat(*edit.AlbumPeak, 'f', 6, 64)}
	}
	if edit.TrackGain != nil {
		tags[models.TagTrackGain] = []string{strconv.FormatFloat(*edit.TrackGain, 'f', 6, 64)}
	}
	if edit.TrackPeak != nil {
		tags[models.TagTrackPeak] = []string{strconv.FormatFloat(*edit.TrackPeak, 'f', 6, 64)}
	}
	if edit.Isrc != nil {
		if *edit.Isrc == "" {
			utils.GinPrettyError(c, http.StatusBadRequest,
				errors.New("isrc field empty"))
			return
		} else {
			tags[models.TagISRC] = []string{*edit.Isrc}
			song.Isrc = *edit.Isrc
		}
	}

	extension := filepath.Ext(song.Path)

	newSongPath := filepath.Join(strings.ReplaceAll(artist, "/", "_"), strings.ReplaceAll(album, "/", "_"),
		fmt.Sprintf("%s - %s%s", strings.ReplaceAll(trackNumber, "/", "_"), strings.ReplaceAll(title, "/", "_"), extension))
	if song.Path != newSongPath {
		root, err := os.OpenRoot(userPath)
		if err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError,
				fmt.Errorf("os.OpenRoot: %w", err))
			return
		}
		defer root.Close()

		oldSongPath := song.Path
		if err := utils.RootDuplicateFile(root, oldSongPath, newSongPath); err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}

		if err := metadata.WriteTags(filepath.Join(userPath, newSongPath), tags, false); err != nil {
			_ = root.Remove(newSongPath)
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		} else {
			_ = root.Remove(oldSongPath)
		}
		if cover != nil {
			if err := metadata.WriteCover(filepath.Join(userPath, newSongPath), cover); err != nil {
				_ = root.Remove(newSongPath)
				utils.GinPrettyError(c, http.StatusInternalServerError, err)
				return
			} else {
				_ = root.Remove(oldSongPath)
			}
		}

		song.Path = newSongPath
	} else {
		if err := metadata.WriteTags(filepath.Join(userPath, newSongPath), tags, false); err != nil {
			utils.GinPrettyError(c, http.StatusInternalServerError, err)
			return
		}
		if cover != nil {
			if err := metadata.WriteCover(filepath.Join(userPath, newSongPath), cover); err != nil {
				utils.GinPrettyError(c, http.StatusInternalServerError, err)
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
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	result, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("strconv.ParseUint: %w", err))
		return
	}
	id := uint(result)

	img, err := services.GetLibrarySongCover(userId, id)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "image/jpg", img)
}

func ListSong(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	var limit int
	if result, err := strconv.Atoi(c.Query("limit")); err != nil {
		limit = 10
	} else {
		limit = result
	}

	var offset int
	if result, err := strconv.Atoi(c.Query("offset")); err != nil {
		offset = 0
	} else {
		offset = result
	}

	q := c.Query("q")
	q = strings.ReplaceAll(q, "_", "\\_")
	q = strings.ReplaceAll(q, "/", "\\_")

	total, err := repository.CountSongByUserID(userId, q)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	songs, err := repository.ListSongByUserID(userId, q, limit, offset)
	if err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
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
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	result, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest,
			fmt.Errorf("strconv.ParseUint: %w", err))
		return
	}
	id := uint(result)

	if err := services.DeleteLibrarySong(userId, id); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func SyncLibrary(c *gin.Context) {
	userId, err := utils.GetFromContext[uint](c, "userId")
	if err != nil {
		utils.GinPrettyError(c, http.StatusBadRequest, err)
		return
	}

	if err := services.SyncUserLibrary(userId); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
