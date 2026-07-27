package response

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type AudioQuality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type SongInfoArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SongInfoAlbum struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type SongInfo struct {
	Provider        string           `json:"provider"`
	ID              string           `json:"id"`
	Title           string           `json:"title"`
	Duration        uint             `json:"duration"`
	ReplayGain      float64          `json:"replayGain"`
	Peak            float64          `json:"peak"`
	AlbumReplayGain float64          `json:"albumReplayGain"`
	AlbumPeak       float64          `json:"albumPeak"`
	ReleaseDate     time.Time        `json:"releaseDate"`
	TrackNumber     uint             `json:"trackNumber"`
	VolumeNumber    uint             `json:"volumeNumber"`
	AudioQuality    AudioQuality     `json:"audioQuality"`
	Explicit        bool             `json:"explicit"`
	Popularity      uint             `json:"popularity"`
	Isrc            string           `json:"isrc"`
	Artists         []SongInfoArtist `json:"artists"`
	Album           SongInfoAlbum    `json:"album"`
}

func SongInfoToResponse(songInfo model.EnrichedSongInfo) SongInfo {
	artists := []SongInfoArtist{}
	for _, a := range songInfo.Artists {
		artists = append(artists, SongInfoArtist{
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
		Album: SongInfoAlbum{
			ID:       songInfo.Album.ID,
			Title:    songInfo.Album.Title,
			CoverUrl: songInfo.Album.CoverUrl,
		},
	}
}

type AlbumInfoArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AlbumInfo struct {
	Provider      string            `json:"provider"`
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Duration      uint              `json:"duration"`
	ReleaseDate   time.Time         `json:"releaseDate"`
	NumberTracks  uint              `json:"numberTracks"`
	NumberVolumes uint              `json:"numberVolumes"`
	CoverUrl      string            `json:"coverUrl"`
	AudioQuality  AudioQuality      `json:"audioQuality"`
	Explicit      bool              `json:"explicit"`
	Artists       []AlbumInfoArtist `json:"artists"`
}

func AlbumInfoToResponse(albumInfo model.EnrichedAlbumInfo) AlbumInfo {
	artists := []AlbumInfoArtist{}
	for _, a := range albumInfo.Artists {
		artists = append(artists, AlbumInfoArtist{
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
		CoverUrl:      albumInfo.CoverUrl,
		AudioQuality: AudioQuality{
			Name:  albumInfo.AudioQuality.Name,
			Color: albumInfo.AudioQuality.Color,
		},
		Explicit: albumInfo.Explicit,
		Artists:  artists,
	}
}

type AlbumSongArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AlbumSong struct {
	Provider     string            `json:"provider"`
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	TrackNumber  uint              `json:"trackNumber"`
	VolumeNumber uint              `json:"volumeNumber"`
	AudioQuality AudioQuality      `json:"audioQuality"`
	Explicit     bool              `json:"explicit"`
	Isrc         string            `json:"isrc"`
	Artists      []AlbumSongArtist `json:"artists"`
}

type AlbumSongs struct {
	Provider string      `json:"provider"`
	Songs    []AlbumSong `json:"songs"`
}

func AlbumSongsToResponse(albumSongs model.EnrichedAlbumSongs) AlbumSongs {
	songs := []AlbumSong{}
	for _, s := range albumSongs.Songs {
		songArtists := []AlbumSongArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, AlbumSongArtist{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		songs = append(songs, AlbumSong{
			Provider:     s.Provider,
			ID:           s.ID,
			Title:        s.Title,
			Duration:     s.Duration,
			TrackNumber:  s.TrackNumber,
			VolumeNumber: s.VolumeNumber,
			AudioQuality: AudioQuality{
				Name:  s.AudioQuality.Name,
				Color: s.AudioQuality.Color,
			},
			Explicit: s.Explicit,
			Isrc:     s.Isrc,
			Artists:  songArtists,
		})
	}

	return AlbumSongs{
		Provider: albumSongs.Provider,
		Songs:    songs,
	}
}

type ArtistInfo struct {
	Provider   string `json:"provider"`
	Followed   string `json:"followed"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
}

func ArtistInfoToResponse(artistInfo model.EnrichedArtistInfo) ArtistInfo {
	return ArtistInfo{
		Provider:   artistInfo.Provider,
		Followed:   artistInfo.Followed.String(),
		ID:         artistInfo.ID,
		Name:       artistInfo.Name,
		PictureUrl: artistInfo.PictureUrl,
	}
}

type ArtistAlbumArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ArtistAlbum struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Duration     uint                `json:"duration"`
	ReleaseDate  time.Time           `json:"releaseDate"`
	CoverUrl     string              `json:"coverUrl"`
	AudioQuality AudioQuality        `json:"audioQuality"`
	Explicit     bool                `json:"explicit"`
	Artists      []ArtistAlbumArtist `json:"artists"`
}

type ArtistAlbums struct {
	Provider string        `json:"provider"`
	Albums   []ArtistAlbum `json:"albums"`
	Ep       []ArtistAlbum `json:"ep"`
	Singles  []ArtistAlbum `json:"singles"`
}

func ArtistAlbumsToResponse(artistAlbums model.EnrichedArtistAlbums) ArtistAlbums {
	albums := []ArtistAlbum{}
	for _, album := range artistAlbums.Albums {
		albumArtists := []ArtistAlbumArtist{}
		for _, a := range album.Artists {
			albumArtists = append(albumArtists, ArtistAlbumArtist{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		albums = append(albums, ArtistAlbum{
			ID:          album.ID,
			Title:       album.Title,
			Duration:    album.Duration,
			ReleaseDate: album.ReleaseDate,
			CoverUrl:    album.CoverUrl,
			AudioQuality: AudioQuality{
				Name:  album.AudioQuality.Name,
				Color: album.AudioQuality.Color,
			},
			Explicit: album.Explicit,
			Artists:  albumArtists,
		})
	}

	eps := []ArtistAlbum{}
	for _, ep := range artistAlbums.Ep {
		epArtists := []ArtistAlbumArtist{}
		for _, a := range ep.Artists {
			epArtists = append(epArtists, ArtistAlbumArtist{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		eps = append(eps, ArtistAlbum{
			ID:          ep.ID,
			Title:       ep.Title,
			Duration:    ep.Duration,
			ReleaseDate: ep.ReleaseDate,
			CoverUrl:    ep.CoverUrl,
			AudioQuality: AudioQuality{
				Name:  ep.AudioQuality.Name,
				Color: ep.AudioQuality.Color,
			},
			Explicit: ep.Explicit,
			Artists:  epArtists,
		})
	}

	singles := []ArtistAlbum{}
	for _, single := range artistAlbums.Singles {
		singleArtists := []ArtistAlbumArtist{}
		for _, a := range single.Artists {
			singleArtists = append(singleArtists, ArtistAlbumArtist{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		singles = append(singles, ArtistAlbum{
			ID:          single.ID,
			Title:       single.Title,
			Duration:    single.Duration,
			ReleaseDate: single.ReleaseDate,
			CoverUrl:    single.CoverUrl,
			AudioQuality: AudioQuality{
				Name:  single.AudioQuality.Name,
				Color: single.AudioQuality.Color,
			},
			Explicit: single.Explicit,
			Artists:  singleArtists,
		})
	}

	return ArtistAlbums{
		Provider: artistAlbums.Provider,
		Albums:   albums,
		Ep:       eps,
		Singles:  singles,
	}
}

type PlaylistInfo struct {
	Provider       string    `json:"provider"`
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Duration       uint      `json:"duration"`
	LastUpdated    time.Time `json:"lastUpdated"`
	NumberOfTracks uint      `json:"numberOfTracks"`
	CoverURL       string    `json:"coverUrl"`
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

type PlaylistSongArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PlaylistSong struct {
	ID           string               `json:"id"`
	Title        string               `json:"title"`
	Duration     uint                 `json:"duration"`
	AudioQuality AudioQuality         `json:"audioQuality"`
	Explicit     bool                 `json:"explicit"`
	Isrc         string               `json:"isrc"`
	Artists      []PlaylistSongArtist `json:"artists"`
}

type PlaylistSongs struct {
	Provider string         `json:"provider"`
	Songs    []PlaylistSong `json:"songs"`
}

func PlaylistSongsToResponse(playlistSongs model.EnrichedPlaylistSongs) PlaylistSongs {
	songs := []PlaylistSong{}
	for _, s := range playlistSongs.Songs {
		songArtists := []PlaylistSongArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, PlaylistSongArtist{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		songs = append(songs, PlaylistSong{
			ID:       s.ID,
			Title:    s.Title,
			Duration: s.Duration,
			AudioQuality: AudioQuality{
				Name:  s.AudioQuality.Name,
				Color: s.AudioQuality.Color,
			},
			Explicit: s.Explicit,
			Isrc:     s.Isrc,
			Artists:  songArtists,
		})
	}
	return PlaylistSongs{
		Provider: playlistSongs.Provider,
		Songs:    songs,
	}
}

type SearchItem struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type SearchSongArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchSongAlbum struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type SearchSong struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Duration     uint               `json:"duration"`
	AudioQuality AudioQuality       `json:"audioQuality"`
	Popularity   uint               `json:"popularity"`
	Explicit     bool               `json:"explicit"`
	Isrc         string             `json:"isrc"`
	Artists      []SearchSongArtist `json:"artists"`
	Album        SearchSongAlbum    `json:"album"`
}

type SearchAlbumArtist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchAlbum struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	Duration     uint                `json:"duration"`
	CoverUrl     string              `json:"coverUrl"`
	AudioQuality AudioQuality        `json:"audioQuality"`
	Popularity   uint                `json:"popularity"`
	Explicit     bool                `json:"explicit"`
	Artists      []SearchAlbumArtist `json:"artists"`
}

type SearchArtist struct {
	Followed   string `json:"followed"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}

type SearchPlaylist struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Duration   uint   `json:"duration"`
	CoverURL   string `json:"coverUrl"`
	Popularity uint   `json:"popularity"`
}

type SearchProviderResult struct {
	Songs     []SearchSong     `json:"songs"`
	Albums    []SearchAlbum    `json:"albums"`
	Artists   []SearchArtist   `json:"artists"`
	Playlists []SearchPlaylist `json:"playlists"`
}

type SearchResult struct {
	Item           SearchItem
	ProviderResult map[string]SearchProviderResult
}

func SearchProviderResultToResponse(search model.EnrichedSearch) SearchProviderResult {
	songs := []SearchSong{}
	for _, song := range search.Songs {
		songArtists := []SearchSongArtist{}
		for _, artist := range song.Artists {
			songArtists = append(songArtists, SearchSongArtist{
				ID:   artist.ID,
				Name: artist.Name,
			})
		}

		songs = append(songs, SearchSong{
			ID:           song.ID,
			Title:        song.Title,
			Duration:     song.Duration,
			AudioQuality: AudioQuality{Name: song.AudioQuality.Name, Color: song.AudioQuality.Color},
			Popularity:   song.Popularity,
			Explicit:     song.Explicit,
			Isrc:         song.Isrc,
			Artists:      songArtists,
			Album: SearchSongAlbum{
				ID:       song.Album.ID,
				Title:    song.Album.Title,
				CoverUrl: song.Album.CoverUrl,
			},
		})
	}

	albums := []SearchAlbum{}
	for _, album := range search.Albums {
		albumArtists := []SearchAlbumArtist{}
		for _, artist := range album.Artists {
			albumArtists = append(albumArtists, SearchAlbumArtist{
				ID:   artist.ID,
				Name: artist.Name,
			})
		}

		albums = append(albums, SearchAlbum{
			ID:           album.ID,
			Title:        album.Title,
			Duration:     album.Duration,
			CoverUrl:     album.CoverUrl,
			AudioQuality: AudioQuality{Name: album.AudioQuality.Name, Color: album.AudioQuality.Color},
			Popularity:   album.Popularity,
			Explicit:     album.Explicit,
			Artists:      albumArtists,
		})
	}

	artists := []SearchArtist{}
	for _, artist := range search.Artists {
		artists = append(artists, SearchArtist{
			Followed:   artist.Followed.String(),
			ID:         artist.ID,
			Name:       artist.Name,
			PictureUrl: artist.PictureUrl,
			Popularity: artist.Popularity,
		})
	}

	playlists := []SearchPlaylist{}
	for _, playlist := range search.Playlists {
		playlists = append(playlists, SearchPlaylist{
			ID:         playlist.ID,
			Title:      playlist.Title,
			Duration:   playlist.Duration,
			CoverURL:   playlist.CoverURL,
			Popularity: playlist.Popularity,
		})
	}

	return SearchProviderResult{
		Songs:     songs,
		Albums:    albums,
		Artists:   artists,
		Playlists: playlists,
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
