package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/setup/config"
	"github.com/google/uuid"
)

type SongService struct {
	cfg  config.DownloadConfig
	repo model.SongRepository
}

func NewSongService(cfg config.DownloadConfig, repo model.SongRepository) SongService {
	return SongService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *SongService) AddSong(ctx context.Context, user model.User, songInfo model.SongInfo, filePath string, fileUpdatedAt time.Time) (model.Song, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return model.Song{}, fmt.Errorf("create id for new song for %s: %w", user.ID.String(), err)
	}

	song := model.Song{
		ID:        id,
		UserID:    user.ID,
		ISRC:      songInfo.ISRC,
		Path:      filePath,
		UpdatedAt: fileUpdatedAt,
	}

	song, err = s.repo.CreateSong(ctx, song)
	if err != nil {
		return model.Song{}, fmt.Errorf(": %w", err)
	}

	return song, nil
}

func (s *SongService) Sync(ctx context.Context, user model.User) error {

	return nil
}

func (s *SongService) GetSong(ctx context.Context, user model.User, songInfo model.SongInfo) (model.Song, error) {
	song, err := s.repo.GetSongByFilter(ctx, model.SongFilter{UserID: &user.ID, ISRC: &songInfo.ISRC})
	if err != nil {
		return model.Song{}, fmt.Errorf(": %w", err)
	}

	return song, nil
}

func (s *SongService) ListSong(ctx context.Context, user model.User) ([]model.Song, error) {
	songs, err := s.repo.ListSongsByFilter(ctx, model.SongFilter{UserID: &user.ID})
	if err != nil {
		return []model.Song{}, fmt.Errorf(": %w", err)
	}

	return songs, nil
}

func (s *SongService) RemoveSong(ctx context.Context, user model.User, id uuid.UUID) error {
	err := s.repo.DeleteSongByUserID(ctx, id, user.ID)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	return nil
}

func (s *SongService) RemoveSongByISRC(ctx context.Context, user model.User, isrc string) error {
	err := s.repo.DeleteSongByUserIDByISRC(ctx, isrc, user.ID)
	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	return nil
}
