package models

import (
	"strconv"
)

func (req RequestUploadSong) ToTags() map[string][]string {
	tags := make(map[string][]string)

	if req.Title != nil {
		tags[TagTitle] = []string{*req.Title}
	}

	if req.Album != nil {
		tags[TagAlbum] = []string{*req.Album}
	}

	if req.AlbumArtists != nil && len(*req.AlbumArtists) > 0 {
		tags[TagAlbumArtists] = *req.AlbumArtists
	}

	if req.Artists != nil && len(*req.Artists) > 0 {
		tags[TagArtists] = *req.Artists
	}

	if req.TrackNumber != nil {
		tags[TagTrackNumber] = []string{strconv.FormatUint(uint64(*req.TrackNumber), 10)}
	}

	if req.VolumeNumber != nil {
		tags[TagVolumeNumber] = []string{strconv.FormatUint(uint64(*req.VolumeNumber), 10)}
	}

	if req.ReleaseDate != nil {
		tags[TagReleaseDate] = []string{*req.ReleaseDate}
	}

	if req.Explicit != nil {
		tags[TagExplicit] = []string{strconv.FormatBool(*req.Explicit)}
	}

	if req.AlbumGain != nil {
		tags[TagAlbumGain] = []string{strconv.FormatFloat(*req.AlbumGain, 'f', -1, 64)}
	}

	if req.AlbumPeak != nil {
		tags[TagAlbumPeak] = []string{strconv.FormatFloat(*req.AlbumPeak, 'f', -1, 64)}
	}

	if req.TrackGain != nil {
		tags[TagTrackGain] = []string{strconv.FormatFloat(*req.TrackGain, 'f', -1, 64)}
	}

	if req.TrackPeak != nil {
		tags[TagTrackPeak] = []string{strconv.FormatFloat(*req.TrackPeak, 'f', -1, 64)}
	}

	if req.Isrc != nil {
		tags[TagISRC] = []string{*req.Isrc}
	}

	return tags
}

func (req RequestEditSong) ToTags() map[string][]string {
	tags := make(map[string][]string)

	if req.Title != nil {
		tags[TagTitle] = []string{*req.Title}
	}

	if req.Album != nil {
		tags[TagAlbum] = []string{*req.Album}
	}

	if req.AlbumArtists != nil && len(*req.AlbumArtists) > 0 {
		tags[TagAlbumArtists] = *req.AlbumArtists
	}

	if req.Artists != nil && len(*req.Artists) > 0 {
		tags[TagArtists] = *req.Artists
	}

	if req.TrackNumber != nil {
		tags[TagTrackNumber] = []string{strconv.FormatUint(uint64(*req.TrackNumber), 10)}
	}

	if req.VolumeNumber != nil {
		tags[TagVolumeNumber] = []string{strconv.FormatUint(uint64(*req.VolumeNumber), 10)}
	}

	if req.ReleaseDate != nil {
		tags[TagReleaseDate] = []string{*req.ReleaseDate}
	}

	if req.Explicit != nil {
		tags[TagExplicit] = []string{strconv.FormatBool(*req.Explicit)}
	}

	if req.AlbumGain != nil {
		tags[TagAlbumGain] = []string{strconv.FormatFloat(*req.AlbumGain, 'f', -1, 64)}
	}

	if req.AlbumPeak != nil {
		tags[TagAlbumPeak] = []string{strconv.FormatFloat(*req.AlbumPeak, 'f', -1, 64)}
	}

	if req.TrackGain != nil {
		tags[TagTrackGain] = []string{strconv.FormatFloat(*req.TrackGain, 'f', -1, 64)}
	}

	if req.TrackPeak != nil {
		tags[TagTrackPeak] = []string{strconv.FormatFloat(*req.TrackPeak, 'f', -1, 64)}
	}

	if req.Isrc != nil {
		tags[TagISRC] = []string{*req.Isrc}
	}

	return tags
}
