package services

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils/metadata"
	"go.senan.xyz/taglib"
)

func UploadLibrarySong(ctx context.Context, userId uint, reader io.Reader, extension string, info models.MetadataInfo, cover io.Reader) error {
	user, err := repository.GetUserByID(userId)
	if err != nil {
		return fmt.Errorf("UploadLibrarySong: %w", err)
	}

	root, err := os.OpenRoot(config.LIBRARY_PATH)
	if err != nil {
		return fmt.Errorf("UploadLibrarySong: os.OpenRoot: %w", err)
	}
	defer root.Close()

	if err := root.Mkdir(user.Username, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("UploadLibrarySong: root.Mkdir: %w", err)
	}

	rootUser, err := root.OpenRoot(user.Username)
	if err != nil {
		return fmt.Errorf("UploadLibrarySong: root.OpenRoot: %w", err)
	}
	defer rootUser.Close()

	filename := filepath.Join(strings.ReplaceAll(info.AlbumArtists[0], "/", "_"), strings.ReplaceAll(info.Album, "/", "_"),
		fmt.Sprintf("%s - %s%s", strings.ReplaceAll(info.TrackNumber, "/", "_"), strings.ReplaceAll(info.Title, "/", "_"), extension))
	dirFile := filepath.Dir(filename)

	if err := rootUser.MkdirAll(dirFile, 0755); err != nil {
		return fmt.Errorf("UploadLibrarySong: rootUser.MkdirAll: %w", err)
	}

	file, err := rootUser.Create(filename)
	if err != nil {
		return fmt.Errorf("UploadLibrarySong: rootUser.Create: %w", err)
	}

	if _, err := io.Copy(file, reader); err != nil {
		_ = file.Close()
		if removeErr := rootUser.Remove(filename); removeErr != nil {
			return fmt.Errorf("UploadLibrarySong: io.Copy: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("UploadLibrarySong: io.Copy: %w", err)
		}
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("UploadLibrarySong: file.Close: %w", err)
	}

	path := filepath.Join(root.Name(), rootUser.Name(), filename)

	tags := metadata.MetadataInfoToTags(info)

	if err := metadata.WriteTags(path, tags, false); err != nil {
		if removeErr := rootUser.Remove(filename); removeErr != nil {
			return fmt.Errorf("saveSong: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("saveSong: %w", err)
		}
	}

	if cover != nil {
		if err := metadata.WriteCover(path, cover); err != nil {
			if removeErr := root.Remove(filename); removeErr != nil {
				return fmt.Errorf("saveSong: %w: %w", err, removeErr)
			} else {
				return fmt.Errorf("saveSong: %w", err)
			}
		}
	}

	_ = repository.AddSong(models.Song{UserId: userId, Path: filename, Isrc: info.Isrc})

	return nil
}

func GetLibrarySong(info models.Song) (models.ResponseSong, error) {
	song := models.ResponseSong{
		ID:           info.ID,
		Isrc:         info.Isrc,
		Title:        "Unknown Title",
		Album:        "Unknown Album",
		Artists:      []string{"Unknown Artist"},
		AlbumArtists: []string{"Unknown Artist"},
		ReleaseDate:  "00-00-0000",
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

	if title, ok := tags[metadata.Title]; ok && len(title) > 0 {
		song.Title = title[0]
	}

	if releaseDate, ok := tags[metadata.ReleaseDate]; ok && len(releaseDate) > 0 {
		song.ReleaseDate = releaseDate[0]
	}

	if trackNumber, ok := tags[metadata.TrackNumber]; ok && len(trackNumber) > 0 {
		if trackNumber, err := strconv.ParseUint(trackNumber[0], 10, 0); err == nil {
			song.TrackNumber = uint(trackNumber)
		}
	}

	if volumeNumber, ok := tags[metadata.VolumeNumber]; ok && len(volumeNumber) > 0 {
		if volumeNumber, err := strconv.ParseUint(volumeNumber[0], 10, 0); err == nil {
			song.VolumeNumber = uint(volumeNumber)
		}
	}

	if explicit, ok := tags[metadata.Explicit]; ok && len(explicit) > 0 {
		if explicit, err := strconv.ParseBool(explicit[0]); err == nil {
			song.Explicit = explicit
		}
	}

	if album, ok := tags[metadata.Album]; ok && len(album) > 0 {
		song.Album = album[0]
	}

	if artists, ok := tags[metadata.Artists]; ok && len(artists) > 0 {
		song.Artists = artists
	}

	if artists, ok := tags[metadata.AlbumArtists]; ok && len(artists) > 0 {
		song.AlbumArtists = artists
	}

	if albumGain, ok := tags[metadata.AlbumGain]; ok && len(albumGain) > 0 {
		if albumGain, err := strconv.ParseFloat(albumGain[0], 64); err == nil {
			song.AlbumGain = albumGain
		}
	}

	if albumPeak, ok := tags[metadata.AlbumPeak]; ok && len(albumPeak) > 0 {
		if albumPeak, err := strconv.ParseFloat(albumPeak[0], 64); err == nil {
			song.AlbumPeak = albumPeak
		}
	}

	if trackGain, ok := tags[metadata.TrackGain]; ok && len(trackGain) > 0 {
		if trackGain, err := strconv.ParseFloat(trackGain[0], 64); err == nil {
			song.TrackGain = trackGain
		}
	}

	if trackPeak, ok := tags[metadata.TrackPeak]; ok && len(trackPeak) > 0 {
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
