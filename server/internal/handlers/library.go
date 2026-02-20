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
	"strconv"
	"strings"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			log.Println(fmt.Errorf("tmpFile.Close: %w", err))
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			log.Println(fmt.Errorf("tmpFile.Close: %w", err))
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

	if err := repository.AddSong(models.Song{UserId: userId, Path: path}); err != nil {
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
		log.Println(fmt.Errorf("newFile.Sync: %w", err))
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
	copyFile.Close()
	defer func() {
		if copyFile != nil {
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
	defer func() {
		if copyFile != nil {
			rootUser.RemoveAll(filepath.Dir(path))
		}
	}()

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := database.DB.Model(&models.Song{}).
			Where("id = ? AND user_id = ?", song.ID, userId).
			Updates(models.Song{Path: path}).Error; err != nil {
			return err
		}

		if song.Path == path {
			if err := utils.RenameForce(copyFile.Name(), filepath.Join(userPath, path)); err != nil {
				return err
			}
		} else {
			if err := utils.RenameSoft(copyFile.Name(), filepath.Join(userPath, path)); err != nil {
				return err
			}
			if err := os.Remove(originalFile.Name()); err != nil {
				log.Println(fmt.Errorf("os.Remove: %w", err))
			}
		}

		return nil
	}); err != nil {
		utils.GinPrettyError(c, http.StatusInternalServerError, err)
		return
	}

	copyFile = nil
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
