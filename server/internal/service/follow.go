package service

import (
	"context"
	"fmt"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type FollowService struct {
	plugin *PluginService
	repo   model.FollowRepository
}

func NewFollowService(plugin *PluginService, repo model.FollowRepository) FollowService {
	return FollowService{
		plugin: plugin,
		repo:   repo,
	}
}

func (s *FollowService) CreateFollow(ctx context.Context, follow model.Follow) (model.Follow, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return model.Follow{}, fmt.Errorf("create id for new follow: %w", err)
	}
	follow.ID = id

	artist, err := s.plugin.GetArtistInfo(ctx, follow.UserID, follow.Provider, follow.ArtistID)
	if err != nil {
		return model.Follow{}, fmt.Errorf("get artist info for new follow: %w", err)
	}
	follow.ArtistName = artist.Name
	follow.ArtistPictureURL = artist.PictureURL

	follow, err = s.repo.CreateFollow(ctx, follow)
	if err != nil {
		return model.Follow{}, fmt.Errorf("create new follow: %w", err)
	}

	return follow, nil

}

func (s *FollowService) ListFollowsByFilter(ctx context.Context, filter model.FollowFilter) ([]model.Follow, error) {
	follows, err := s.repo.ListFollowsByFilter(ctx, filter)
	if err != nil {
		return []model.Follow{}, fmt.Errorf("list follows for filter %#v: %w", filter, err)
	}

	return follows, nil
}

func (s *FollowService) DeleteFollow(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteFollow(ctx, id); err != nil {
		return fmt.Errorf("delete follow %s: %w", id.String(), err)
	}

	return nil
}

func (s *FollowService) DeleteFollowByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := s.repo.DeleteFollowByUserID(ctx, id, userID); err != nil {
		return fmt.Errorf("delete follow %s for user %s: %w", id.String(), userID.String(), err)
	}

	return nil
}
