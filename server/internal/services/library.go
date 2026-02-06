package services

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"go.senan.xyz/taglib"
)

func GetLibrarySong(info models.Song) (models.ResponseSong, error) {
	song := models.ResponseSong{
		ID:   info.ID,
		Isrc: info.Isrc,
	}

	path, err := utils.GetUserPath(info.UserId)
	if err != nil {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", err)
	}
	path = filepath.Join(path, info.Path)

	properties, err := taglib.ReadProperties(path)
	if err != nil {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", err)
	}

	song.Duration = uint(properties.Length.Seconds())

	tags, err := taglib.ReadTags(path)
	if err != nil {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", err)
	}

	if title, ok := tags[taglib.Title]; !ok && len(title) == 0 {
		song.Title = "Unknown Title"
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag title not found"))
	} else {
		song.Title = title[0]
	}

	if releaseDate, ok := tags[taglib.ReleaseDate]; !ok && len(releaseDate) == 0 {
		song.ReleaseDate = "00-00-0000"
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag releaseDate not found"))
	} else {
		song.ReleaseDate = releaseDate[0]
	}

	if trackNumber, ok := tags[taglib.TrackNumber]; !ok && len(trackNumber) == 0 {
		song.TrackNumber = 0
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag trackNumber not found"))
	} else {
		if trackNumber, err := strconv.ParseUint(trackNumber[0], 10, 0); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag trackNumber not valid"))
		} else {
			song.TrackNumber = uint(trackNumber)
		}
	}

	if volumeNumber, ok := tags[taglib.DiscNumber]; !ok && len(volumeNumber) == 0 {
		song.VolumeNumber = 0
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag volumeNumber not found"))
	} else {
		if volumeNumber, err := strconv.ParseUint(volumeNumber[0], 10, 0); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag volumeNumber not valid"))
		} else {
			song.VolumeNumber = uint(volumeNumber)
		}
	}

	if explicit, ok := tags["ITUNESADVISORY"]; !ok && len(explicit) == 0 {
		song.Explicit = false
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag explicit not found"))
	} else {
		if explicit, err := strconv.ParseBool(explicit[0]); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag explicit not valid"))
		} else {
			song.Explicit = explicit
		}
	}

	if album, ok := tags[taglib.Album]; !ok && len(album) == 0 {
		song.Album = "Unknown Album"
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag album not found"))
	} else {
		song.Album = album[0]
	}

	if artists, ok := tags[taglib.Artists]; !ok && len(artists) == 0 {
		song.Artists = []string{"Unknown Artist"}
		// return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag artists not found"))
	} else {
		song.Artists = artists
	}

	return song, nil
}

func DeleteLibrarySong(userId uint, id uint) error {
	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		return fmt.Errorf("services.DeleteLibrarySong: %w", err)
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		return fmt.Errorf("services.DeleteLibrarySong: %w", err)
	}

	path := filepath.Join(userPath, song.Path)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("services.DeleteLibrarySong: %w", err)
	}

	dir := filepath.Dir(path)
	root, _ := filepath.Abs(userPath)
	for {
		dirAbs, _ := filepath.Abs(dir)
		if dirAbs == root || dirAbs == "/" {
			break
		}
		if err := os.Remove(dirAbs); err != nil {
			break
		}
		dir = filepath.Dir(dir)
	}

	if err := repository.DeleteSongByUserID(userId, id); err != nil {
		return fmt.Errorf("services.DeleteLibrarySong: %w", err)
	}
	return nil
}

func SyncUserLibrary(userId uint) error {
	tmpDbList, err := repository.ListSongByUserID(userId, -1, 0)
	if err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}

	addList := make(map[string]string)
	type updateItem struct {
		id   uint
		path string
	}
	updateList := make(map[string]updateItem)
	type dbItem struct {
		id   uint
		path string
	}
	deleteList := make(map[string]dbItem)

	diskList := make(map[string]string)

	dbList := make(map[string]dbItem)
	for _, item := range tmpDbList {
		dbList[item.Isrc] = dbItem{
			item.ID,
			item.Path,
		}
		deleteList[item.Isrc] = dbItem{
			item.ID,
			item.Path,
		}
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}

	if err := filepath.WalkDir(userPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		tags, err := taglib.ReadTags(path)
		if err != nil {
			log.Println("services.SyncUserLibrary: filepath.WalkDir", err)
			return nil
		}

		isrc, ok := tags["ISRC"]
		if !ok || len(isrc) == 0 {
			log.Println("services.SyncUserLibrary: filepath.WalkDir: tag ISRC not found")
			return nil
		}

		if _, exists := diskList[isrc[0]]; exists {
			log.Printf("duplicate ISRC on disk: %s", isrc)
			return nil
		}

		path = filepath.Clean(path)
		relPath, err := filepath.Rel(userPath, path)
		if err != nil {
			return err
		}

		diskList[isrc[0]] = relPath
		addList[isrc[0]] = relPath
		return nil
	}); err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}

	for isrc, item := range dbList {
		if path, ok := diskList[isrc]; ok {
			delete(addList, isrc)
			delete(deleteList, isrc)
			if path != item.path {
				updateList[isrc] = updateItem{
					id:   item.id,
					path: path,
				}
			}
		}
	}

	for _, item := range deleteList {
		if err := repository.DeleteSong(item.id); err != nil {
			log.Println("services.SyncUserLibrary:", err)
		}
	}

	for _, item := range updateList {
		if err := repository.UpdateSongByUserID(userId, models.Song{ID: item.id, Path: item.path}); err != nil {
			log.Println("services.SyncUserLibrary:", err)
		}
	}

	for isrc, path := range addList {
		if err := repository.AddSong(models.Song{UserId: userId, Path: path, Isrc: isrc}); err != nil {
			log.Println("services.SyncUserLibrary:", err)
		}
	}

	return nil
}
