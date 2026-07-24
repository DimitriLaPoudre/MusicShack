package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type FetchNewReleasesService struct {
	plugin   *PluginService
	follow   *FollowService
	download *DownloadService
}

type release struct {
	userID   uuid.UUID
	provider string
	albumID  string
}

func NewFetchNewReleasesService(plugin *PluginService, follow *FollowService, download *DownloadService) FetchNewReleasesService {
	return FetchNewReleasesService{
		plugin:   plugin,
		follow:   follow,
		download: download,
	}
}

func (s *FetchNewReleasesService) getArtistNewReleases(ctx context.Context, follow model.Follow, lastFetchDate time.Time) ([]release, error) {
	artist, err := s.plugin.GetArtist(ctx, follow.UserID, follow.Provider, follow.ArtistID)
	if err != nil {
		return []release{}, fmt.Errorf("get artist %s: %w", follow.ArtistID, err)
	}

	var releases []model.EnrichedAlbum
	releases = append(releases, artist.Albums...)
	releases = append(releases, artist.Ep...)
	releases = append(releases, artist.Singles...)

	var newReleases []release
	for _, r := range releases {
		if r.ReleaseDate.After(lastFetchDate) {
			newReleases = append(newReleases, release{
				userID:   follow.UserID,
				provider: artist.Provider,
				albumID:  r.ID,
			})
		}
	}

	return newReleases, nil
}

func (s *FetchNewReleasesService) getNewReleases(ctx context.Context, follows []model.Follow, lastFetchDate time.Time) ([]release, []error) {
	var newReleases []release
	var errList []error
	for _, follow := range follows {
		if tmp, err := s.getArtistNewReleases(ctx, follow, lastFetchDate); err != nil {
			errList = append(errList, fmt.Errorf("get artist new releases: %w", err))
		} else {
			newReleases = append(newReleases, tmp...)
		}
	}

	return newReleases, errList
}

func (s *FetchNewReleasesService) FetchNewReleases(ctx context.Context, lastFetchDate time.Time) []error {
	follows, err := s.follow.ListFollowsByFilter(ctx, model.FollowFilter{})
	if err != nil {
		return []error{err}
	}

	var errList []error
	newReleases, errList := s.getNewReleases(ctx, follows, lastFetchDate)

	for _, release := range newReleases {
		if err := s.download.DownloadAlbum(ctx, release.userID, release.provider, release.albumID); err != nil {
			errList = append(errList, fmt.Errorf("download album %s: %w", release.albumID, err))
		}
	}

	return errList
}
