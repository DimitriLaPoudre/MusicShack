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

func UploadLibrarySong(ctx context.Context, userId uint, reader io.ReadCloser, extension string, info models.MetadataInfo) error {
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
		fmt.Sprintf("%s - %s.%s", strings.ReplaceAll(info.TrackNumber, "/", "_"), strings.ReplaceAll(info.Title, "/", "_"), extension))
	dirFile := filepath.Dir(filename)

	if err := rootUser.MkdirAll(dirFile, 0755); err != nil {
		return fmt.Errorf("UploadLibrarySong: rootUser.MkdirAll: %w", err)
	}

	file, err := rootUser.Create(filename)
	if err != nil {
		return fmt.Errorf("UploadLibrarySong: rootUser.Create: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		filepath := file.Name()
		_ = file.Close()
		if removeErr := os.Remove(filepath); removeErr != nil {
			return fmt.Errorf("saveSong: io.Copy: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("saveSong: io.Copy: %w", err)
		}
	}

	if err := metadata.ApplyMetadata(filepath.Join(root.Name(), rootUser.Name(), filename), info); err != nil {
		filepath := file.Name()
		_ = file.Close()
		if removeErr := os.Remove(filepath); removeErr != nil {
			return fmt.Errorf("saveSong: %w: %w", err, removeErr)
		} else {
			return fmt.Errorf("saveSong: %w", err)
		}
	}

	_ = repository.AddSong(models.Song{UserId: userId, Path: filename, Isrc: info.Isrc})

	return nil
}

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
	} else {
		song.Title = title[0]
	}

	if releaseDate, ok := tags[taglib.ReleaseDate]; !ok && len(releaseDate) == 0 {
		song.ReleaseDate = "00-00-0000"
	} else {
		song.ReleaseDate = releaseDate[0]
	}

	if trackNumber, ok := tags[taglib.TrackNumber]; !ok && len(trackNumber) == 0 {
		song.TrackNumber = 0
	} else {
		if trackNumber, err := strconv.ParseUint(trackNumber[0], 10, 0); err != nil {
			song.TrackNumber = 0
		} else {
			song.TrackNumber = uint(trackNumber)
		}
	}

	if volumeNumber, ok := tags[taglib.DiscNumber]; !ok && len(volumeNumber) == 0 {
		song.VolumeNumber = 0
	} else {
		if volumeNumber, err := strconv.ParseUint(volumeNumber[0], 10, 0); err != nil {
			song.VolumeNumber = 0
		} else {
			song.VolumeNumber = uint(volumeNumber)
		}
	}

	if explicit, ok := tags["ITUNESADVISORY"]; !ok && len(explicit) == 0 {
		song.Explicit = false
	} else {
		if explicit, err := strconv.ParseBool(explicit[0]); err != nil {
			song.Explicit = false
		} else {
			song.Explicit = explicit
		}
	}

	if album, ok := tags[taglib.Album]; !ok && len(album) == 0 {
		song.Album = "Unknown Album"
	} else {
		song.Album = album[0]
	}

	if artists, ok := tags[taglib.Artists]; !ok && len(artists) == 0 {
		song.Artists = []string{"Unknown Artist"}
	} else {
		song.Artists = artists
	}

	if artists, ok := tags["ALBUMARTISTS"]; !ok && len(artists) == 0 {
		song.AlbumArtists = []string{"Unknown Artist"}
	} else {
		song.AlbumArtists = artists
	}

	if albumGain, ok := tags["REPLAYGAIN_ALBUM_GAIN"]; !ok && len(albumGain) == 0 {
		song.AlbumGain = 0
	} else {
		if albumGain, err := strconv.ParseFloat(albumGain[0], 64); err != nil {
			song.AlbumGain = 0
		} else {
			song.AlbumGain = albumGain
		}
	}

	if albumPeak, ok := tags["REPLAYGAIN_ALBUM_PEAK"]; !ok && len(albumPeak) == 0 {
		song.AlbumPeak = 0
	} else {
		if albumPeak, err := strconv.ParseFloat(albumPeak[0], 64); err != nil {
			song.AlbumPeak = 0
		} else {
			song.AlbumPeak = albumPeak
		}
	}

	if trackGain, ok := tags["REPLAYGAIN_TRACK_GAIN"]; !ok && len(trackGain) == 0 {
		song.TrackGain = 0
	} else {
		if trackGain, err := strconv.ParseFloat(trackGain[0], 64); err != nil {
			song.TrackGain = 0
		} else {
			song.TrackGain = trackGain
		}
	}

	if trackPeak, ok := tags["REPLAYGAIN_TRACK_PEAK"]; !ok && len(trackPeak) == 0 {
		song.TrackPeak = 0
	} else {
		if trackPeak, err := strconv.ParseFloat(trackPeak[0], 64); err != nil {
			song.TrackPeak = 0
		} else {
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

	img, err := taglib.ReadImage(path)
	if err != nil {
		return nil, fmt.Errorf("services.GetLibrarySongCover: taglib.ReadImage: %w", err)
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
	tmpDbList, err := repository.ListSongByUserID(userId, "", -1, 0)
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
