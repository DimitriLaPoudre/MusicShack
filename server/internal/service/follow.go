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

func (s *FollowService) Create(ctx context.Context, user model.User, follow model.Follow) (model.Follow, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return model.Follow{}, fmt.Errorf("create id for new follow: %w", err)
	}
	follow.ID = id

	artist, err := s.plugin.GetArtist(ctx, user, follow.Provider, follow.ArtistID)
	follow.ArtistName = artist.Name
	follow.ArtistPictureURL = artist.PictureUrl

	follow, err = s.repo.CreateFollow(ctx, follow)
	if err != nil {
		return model.Follow{}, fmt.Errorf("create new follow: %w", err)
	}

	return follow, nil

}

func (s *FollowService) ListWithFilter(c context.Context, filter model.FollowFilter) ([]model.Follow, error) {
	follows, err := s.repo.ListFollowsByFilter(c, filter)
	if err != nil {
		return []model.Follow{}, fmt.Errorf("list follows for filter %#v: %w", filter, err)
	}

	return follows, nil
}

func (s *FollowService) Delete(c context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteFollow(c, id); err != nil {
		return fmt.Errorf("delete follow %s: %w", id.String(), err)
	}

	return nil
}

func (s *FollowService) DeleteByUserID(c context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := s.repo.DeleteFollowByUserID(c, id, userID); err != nil {
		return fmt.Errorf("delete follow %s for user %s: %w", id.String(), userID.String(), err)
	}

	return nil
}
