package services

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"go.senan.xyz/taglib"
)

func GetLibrarySong(info models.Song) (models.ResponseSong, error) {
	song := models.ResponseSong{
		ID:   info.ID,
		Isrc: info.Isrc,
	}

	properties, err := taglib.ReadProperties(info.Path)
	if err != nil {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", err)
	}

	song.Duration = uint(properties.Length.Seconds())

	tags, err := taglib.ReadTags(info.Path)
	if err != nil {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", err)
	}

	if title, ok := tags[taglib.Title]; !ok && len(title) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag title not found"))
	} else {
		song.Title = title[0]
	}

	if releaseDate, ok := tags[taglib.ReleaseDate]; !ok && len(releaseDate) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag releaseDate not found"))
	} else {
		song.ReleaseDate = releaseDate[0]
	}

	if trackNumber, ok := tags[taglib.TrackNumber]; !ok && len(trackNumber) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag trackNumber not found"))
	} else {
		if trackNumber, err := strconv.ParseUint(trackNumber[0], 10, 0); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag trackNumber not valid"))
		} else {
			song.TrackNumber = uint(trackNumber)
		}
	}

	if volumeNumber, ok := tags[taglib.DiscNumber]; !ok && len(volumeNumber) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag volumeNumber not found"))
	} else {
		if volumeNumber, err := strconv.ParseUint(volumeNumber[0], 10, 0); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag volumeNumber not valid"))
		} else {
			song.VolumeNumber = uint(volumeNumber)
		}
	}

	if explicit, ok := tags["ITUNESADVISORY"]; !ok && len(explicit) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag explicit not found"))
	} else {
		if explicit, err := strconv.ParseBool(explicit[0]); err != nil {
			return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag explicit not valid"))
		} else {
			song.Explicit = explicit
		}
	}

	if album, ok := tags[taglib.Album]; !ok && len(album) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag album not found"))
	} else {
		song.Album = album[0]
	}

	if artists, ok := tags[taglib.Artists]; !ok && len(artists) == 0 {
		return models.ResponseSong{}, fmt.Errorf("services.GetLibrarySong: %w", errors.New("tag artists not found"))
	} else {
		song.Artists = artists
	}

	return song, nil
}

func SyncUserLibrary(userId uint) error {
	user, err := repository.GetUserByID(userId)
	if err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}

	if err := filepath.WalkDir(filepath.Join(config.LIBRARY_PATH, user.Username), func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		tags, err2 := taglib.ReadTags(path)
		if err2 != nil {
			log.Println("services.WalkDir:", err2)
			return nil
		}

		isrc, ok := tags["ISRC"]
		if !ok || len(isrc) == 0 {
			log.Println("services.WalkDir: tag ISRC not found")
			return nil
		}

		if err := repository.SaveSong(models.Song{
			UserId: userId,
			Isrc:   isrc[0],
			Path:   path,
		}); err != nil {
			log.Println("services.WalkDir:", err)
			return nil
		}

		return nil
	}); err != nil {
		return fmt.Errorf("services.SyncUserLibrary: %w", err)
	}

	return nil
}
