package response

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type AudioQuality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Pagination struct {
	Limit              int `json:"limit"`
	Offset             int `json:"offset"`
	TotalNumberOfItems int `json:"total_number_of_items"`
}

// -- SongInfo -- //

type SongInfo struct {
	Provider        string       `json:"provider"`
	ID              string       `json:"id"`
	Title           string       `json:"title"`
	Duration        int          `json:"duration"`
	ReplayGain      float64      `json:"replay_gain"`
	Peak            float64      `json:"peak"`
	AlbumReplayGain float64      `json:"album_replay_gain"`
	AlbumPeak       float64      `json:"album_peak"`
	ReleaseDate     time.Time    `json:"release_date"`
	TrackNumber     int          `json:"track_number"`
	VolumeNumber    int          `json:"volume_number"`
	AudioQuality    AudioQuality `json:"audio_quality"`
	Explicit        bool         `json:"explicit"`
	Popularity      int          `json:"popularity"`
	Isrc            string       `json:"isrc"`
	Artists         []MiniArtist `json:"artists"`
	Album           SongAlbum    `json:"album"`
}

type SongAlbum struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	CoverURL string `json:"cover_url"`
}

func SongInfoToResponse(songInfo model.EnrichedSongInfo) SongInfo {
	artists := []MiniArtist{}
	for _, a := range songInfo.Artists {
		artists = append(artists, MiniArtist{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return SongInfo{
		Provider:        songInfo.Provider,
		ID:              songInfo.ID,
		Title:           songInfo.Title,
		Duration:        songInfo.Duration,
		ReplayGain:      songInfo.ReplayGain,
		Peak:            songInfo.Peak,
		AlbumReplayGain: songInfo.AlbumReplayGain,
		AlbumPeak:       songInfo.AlbumPeak,
		ReleaseDate:     songInfo.ReleaseDate,
		TrackNumber:     songInfo.TrackNumber,
		VolumeNumber:    songInfo.VolumeNumber,
		AudioQuality: AudioQuality{
			Name:  songInfo.AudioQuality.Name,
			Color: songInfo.AudioQuality.Color,
		},
		Explicit:   songInfo.Explicit,
		Popularity: songInfo.Popularity,
		Isrc:       songInfo.Isrc,
		Artists:    artists,
		Album: SongAlbum{
			ID:       songInfo.Album.ID,
			Title:    songInfo.Album.Title,
			CoverURL: songInfo.Album.CoverURL,
		},
	}
}

// -- PagignatedSongs -- //

type PaginatedSongs struct {
	Provider string `json:"provider"`
	Pagination
	Songs []SongInfo `json:"items"`
}

func PaginatedSongsToResponse(paginatedSongs model.EnrichedPaginatedSongs) PaginatedSongs {
	songs := []SongInfo{}
	for _, s := range paginatedSongs.Songs {
		songs = append(songs, SongInfoToResponse(s))
	}

	return PaginatedSongs{
		Provider: paginatedSongs.Provider,
		Songs:    songs,
	}
}

// -- AlbumInfo -- //

type AlbumInfo struct {
	Provider      string       `json:"provider"`
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Duration      int          `json:"duration"`
	ReleaseDate   time.Time    `json:"release_date"`
	NumberTracks  int          `json:"number_tracks"`
	NumberVolumes int          `json:"number_volumes"`
	CoverURL      string       `json:"cover_url"`
	AudioQuality  AudioQuality `json:"audio_quality"`
	Explicit      bool         `json:"explicit"`
	Artists       []MiniArtist `json:"artists"`
}

func AlbumInfoToResponse(albumInfo model.EnrichedAlbumInfo) AlbumInfo {
	artists := []MiniArtist{}
	for _, a := range albumInfo.Artists {
		artists = append(artists, MiniArtist{
			ID:   a.ID,
			Name: a.Name,
		})
	}

	return AlbumInfo{
		Provider:      albumInfo.Provider,
		ID:            albumInfo.ID,
		Title:         albumInfo.Title,
		Duration:      albumInfo.Duration,
		ReleaseDate:   albumInfo.ReleaseDate,
		NumberTracks:  albumInfo.NumberTracks,
		NumberVolumes: albumInfo.NumberVolumes,
		CoverURL:      albumInfo.CoverURL,
		AudioQuality: AudioQuality{
			Name:  albumInfo.AudioQuality.Name,
			Color: albumInfo.AudioQuality.Color,
		},
		Explicit: albumInfo.Explicit,
		Artists:  artists,
	}
}

// -- PaginatedAlbums -- //

type PaginatedAlbums struct {
	Provider string `json:"provider"`
	Pagination
	Albums []AlbumInfo `json:"items"`
}

func PaginatedAlbumsToResponse(paginatedAlbums model.EnrichedPaginatedAlbums) PaginatedAlbums {
	albums := []AlbumInfo{}
	for _, a := range paginatedAlbums.Albums {
		albums = append(albums, AlbumInfoToResponse(a))
	}

	return PaginatedAlbums{
		Provider: paginatedAlbums.Provider,
		Albums:   albums,
	}
}

// -- MiniArtist -- //

type MiniArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// -- ArtistInfo -- //

type ArtistInfo struct {
	Provider   string  `json:"provider"`
	Followed   *string `json:"followed"`
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	PictureURL string  `json:"picture_url"`
}

func ArtistInfoToResponse(artistInfo model.EnrichedArtistInfo) ArtistInfo {
	var followed *string
	if artistInfo.Followed != nil {
		followed_ground := (*artistInfo.Followed).String()
		followed = &followed_ground
	}
	return ArtistInfo{
		Provider:   artistInfo.Provider,
		Followed:   followed,
		ID:         artistInfo.ID,
		Name:       artistInfo.Name,
		PictureURL: artistInfo.PictureURL,
	}
}

// -- PaginatedArtists -- //

type PaginatedArtists struct {
	Provider string `json:"provider"`
	Pagination
	Artists []ArtistInfo `json:"items"`
}

func PaginatedArtistsToResponse(paginatedArtists model.EnrichedPaginatedArtists) PaginatedArtists {
	artists := []ArtistInfo{}
	for _, a := range paginatedArtists.Artists {
		artists = append(artists, ArtistInfoToResponse(a))
	}

	return PaginatedArtists{
		Provider: paginatedArtists.Provider,
		Artists:  artists,
	}
}

// -- ArtistPaginatedAlbums -- //

type ArtistPaginatedAlbums struct {
	Provider string          `json:"provider"`
	Albums   PaginatedAlbums `json:"albums"`
	EPs      PaginatedAlbums `json:"ep"`
	Singles  PaginatedAlbums `json:"singles"`
}

func ArtistPaginatedAlbumsToResponse(artistPaginatedAlbums model.EnrichedArtistPaginatedAlbums) ArtistPaginatedAlbums {
	return ArtistPaginatedAlbums{
		Provider: artistPaginatedAlbums.Provider,
		Albums:   PaginatedAlbumsToResponse(artistPaginatedAlbums.Albums),
		EPs:      PaginatedAlbumsToResponse(artistPaginatedAlbums.EPs),
		Singles:  PaginatedAlbumsToResponse(artistPaginatedAlbums.Singles),
	}
}

// -- PlaylistInfo -- //

type PlaylistInfo struct {
	Provider       string    `json:"provider"`
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Duration       int       `json:"duration"`
	LastUpdated    time.Time `json:"last_updated"`
	NumberOfTracks int       `json:"number_of_tracks"`
	CoverURL       string    `json:"cover_url"`
}

func PlaylistInfoToResponse(playlistInfo model.EnrichedPlaylistInfo) PlaylistInfo {
	return PlaylistInfo{
		Provider:       playlistInfo.Provider,
		ID:             playlistInfo.ID,
		Title:          playlistInfo.Title,
		Description:    playlistInfo.Description,
		Duration:       playlistInfo.Duration,
		LastUpdated:    playlistInfo.LastUpdated,
		NumberOfTracks: playlistInfo.NumberOfTracks,
		CoverURL:       playlistInfo.CoverURL,
	}
}

// -- PaginatedPlaylists -- //
type PaginatedPlaylists struct {
	Provider string `json:"provider"`
	Pagination
	Playlists []PlaylistInfo `json:"items"`
}

func PaginatedPlaylistsToResponse(paginatedPlaylists model.EnrichedPaginatedPlaylists) PaginatedPlaylists {
	playlists := []PlaylistInfo{}
	for _, a := range paginatedPlaylists.Playlists {
		playlists = append(playlists, PlaylistInfoToResponse(a))
	}

	return PaginatedPlaylists{
		Provider:  paginatedPlaylists.Provider,
		Playlists: playlists,
	}
}

// -- Search -- //

type SearchItem struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type SearchProviderResult struct {
	Songs     PaginatedSongs     `json:"songs"`
	Albums    PaginatedAlbums    `json:"albums"`
	Artists   PaginatedArtists   `json:"artists"`
	Playlists PaginatedPlaylists `json:"playlists"`
}

type SearchResult struct {
	Item           SearchItem                      `json:"item"`
	ProviderResult map[string]SearchProviderResult `json:"provider_result"`
}

func SearchProviderResultToResponse(search model.EnrichedSearch) SearchProviderResult {
	return SearchProviderResult{
		Songs:     PaginatedSongsToResponse(search.Songs),
		Albums:    PaginatedAlbumsToResponse(search.Albums),
		Artists:   PaginatedArtistsToResponse(search.Artists),
		Playlists: PaginatedPlaylistsToResponse(search.Playlists),
	}
}

func SearchResultToResponse(searchResult model.SearchResult) SearchResult {
	if searchResult.ProviderResult == nil {
		return SearchResult{
			Item: SearchItem{
				Type: string(searchResult.ItemFound.Type),
				Data: searchResult.ItemFound.Data,
			},
		}
	}

	providerResult := map[string]SearchProviderResult{}
	for provider, result := range searchResult.ProviderResult {
		providerResult[provider] = SearchProviderResultToResponse(result)
	}
	return SearchResult{
		ProviderResult: providerResult,
	}
}
