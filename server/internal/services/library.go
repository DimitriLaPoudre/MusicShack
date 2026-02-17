package services

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
	"go.senan.xyz/taglib"
)

func GetSongPathByTags(tags map[string][]string, extension string) (string, error) {
	var title string
	if value, ok := tags[models.TagTitle]; !ok || len(value) <= 0 || len(value[0]) <= 0 {
		return "", errors.New("title field empty and file don't provide default")
	} else {
		title = strings.ReplaceAll(value[0], "/", "_")
	}
	var album string
	if value, ok := tags[models.TagAlbum]; !ok || len(value) <= 0 || len(value[0]) <= 0 {
		return "", errors.New("album field empty and file don't provide default")
	} else {
		album = strings.ReplaceAll(value[0], "/", "_")
	}
	if value, ok := tags[models.TagArtists]; !ok || len(value) <= 0 || len(value[0]) <= 0 {
		return "", errors.New("artists field empty and file don't provide default")
	}
	var albumArtist string
	if value, ok := tags[models.TagAlbumArtists]; !ok || len(value) <= 0 || len(value[0]) <= 0 {
		return "", errors.New("albumArtists field empty and file don't provide default")
	} else {
		albumArtist = strings.ReplaceAll(value[0], "/", "_")
	}
	var trackNumber string
	if value, ok := tags[models.TagTrackNumber]; !ok || len(value) <= 0 || len(value[0]) <= 0 {
		trackNumber = "0"
	} else {
		trackNumber = strings.ReplaceAll(value[0], "/", "_")
	}

	return filepath.Join(albumArtist, album, fmt.Sprintf("%s - %s%s", trackNumber, title, extension)), nil
}

func GetLibrarySong(info models.Song) (models.ResponseSong, error) {
	song := models.ResponseSong{
		ID:           info.ID,
		Isrc:         info.Isrc,
		Title:        "",
		Album:        "",
		Artists:      []string{},
		AlbumArtists: []string{},
		ReleaseDate:  "",
		TrackNumber:  0,
		VolumeNumber: 0,
		Explicit:     false,
		Duration:     0,
		AlbumGain:    0,
		AlbumPeak:    0,
		TrackGain:    0,
		TrackPeak:    0,
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

	if title, ok := tags[models.TagTitle]; ok && len(title) > 0 {
		song.Title = title[0]
	}

	if releaseDate, ok := tags[models.TagReleaseDate]; ok && len(releaseDate) > 0 {
		song.ReleaseDate = releaseDate[0]
	}

	if trackNumber, ok := tags[models.TagTrackNumber]; ok && len(trackNumber) > 0 {
		if trackNumber, err := strconv.ParseUint(trackNumber[0], 10, 0); err == nil {
			song.TrackNumber = uint(trackNumber)
		}
	}

	if volumeNumber, ok := tags[models.TagVolumeNumber]; ok && len(volumeNumber) > 0 {
		if volumeNumber, err := strconv.ParseUint(volumeNumber[0], 10, 0); err == nil {
			song.VolumeNumber = uint(volumeNumber)
		}
	}

	if explicit, ok := tags[models.TagExplicit]; ok && len(explicit) > 0 {
		if explicit, err := strconv.ParseBool(explicit[0]); err == nil {
			song.Explicit = explicit
		}
	}

	if album, ok := tags[models.TagAlbum]; ok && len(album) > 0 {
		song.Album = album[0]
	}

	if artists, ok := tags[models.TagArtists]; ok && len(artists) > 0 {
		song.Artists = artists
	}

	if artists, ok := tags[models.TagAlbumArtists]; ok && len(artists) > 0 {
		song.AlbumArtists = artists
	}

	if albumGain, ok := tags[models.TagAlbumGain]; ok && len(albumGain) > 0 {
		if albumGain, err := strconv.ParseFloat(albumGain[0], 64); err == nil {
			song.AlbumGain = albumGain
		}
	}

	if albumPeak, ok := tags[models.TagAlbumPeak]; ok && len(albumPeak) > 0 {
		if albumPeak, err := strconv.ParseFloat(albumPeak[0], 64); err == nil {
			song.AlbumPeak = albumPeak
		}
	}

	if trackGain, ok := tags[models.TagTrackGain]; ok && len(trackGain) > 0 {
		if trackGain, err := strconv.ParseFloat(trackGain[0], 64); err == nil {
			song.TrackGain = trackGain
		}
	}

	if trackPeak, ok := tags[models.TagTrackPeak]; ok && len(trackPeak) > 0 {
		if trackPeak, err := strconv.ParseFloat(trackPeak[0], 64); err == nil {
			song.TrackPeak = trackPeak
		}
	}

	return song, nil
}

func GetLibrarySongCover(userId uint, id uint) ([]byte, error) {
	song, err := repository.GetSongByUserID(userId, id)
	if err != nil {
		return nil, fmt.Errorf("services.GetLibrarySongCover: %w", err)
	}

	userPath, err := utils.GetUserPath(userId)
	if err != nil {
		return nil, fmt.Errorf("services.GetLibrarySongCover: %w", err)
	}

	path := filepath.Join(userPath, song.Path)

	img, err := metadata.ReadCover(path)
	if err != nil {
		return nil, fmt.Errorf("services.GetLibrarySongCover: %w", err)
	}

	return img, nil
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

		tags, err := metadata.ReadTags(path)
		if err != nil {
			log.Println("services.SyncUserLibrary: filepath.WalkDir:", path, err)
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

	dbList := make(map[string]dbItem)
	tmpDbList, err := repository.ListSongByUserID(userId, "", -1, 0)
	if err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}
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
