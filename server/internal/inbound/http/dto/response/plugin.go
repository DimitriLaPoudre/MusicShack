package response

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type AudioQuality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type SongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SongAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type Song struct {
	Provider        string       `json:"provider"`
	Downloaded      bool         `json:"downloaded"`
	Id              string       `json:"id"`
	Title           string       `json:"title"`
	Duration        uint         `json:"duration"`
	ReplayGain      float64      `json:"replayGain"`
	Peak            float64      `json:"peak"`
	AlbumReplayGain float64      `json:"albumReplayGain"`
	AlbumPeak       float64      `json:"albumPeak"`
	ReleaseDate     time.Time    `json:"releaseDate"`
	TrackNumber     uint         `json:"trackNumber"`
	VolumeNumber    uint         `json:"volumeNumber"`
	AudioQuality    AudioQuality `json:"audioQuality"`
	Explicit        bool         `json:"explicit"`
	Popularity      uint         `json:"popularity"`
	Isrc            string       `json:"isrc"`
	Artists         []SongArtist `json:"artists"`
	Album           SongAlbum    `json:"album"`
}

func SongToResponse(song model.EnrichedSong) Song {
	artists := []SongArtist{}
	for _, a := range song.Artists {
		artists = append(artists, SongArtist{
			Id:   a.Id,
			Name: a.Name,
		})
	}

	return Song{
		Provider:        song.Provider,
		Downloaded:      song.Downloaded,
		Id:              song.Id,
		Title:           song.Title,
		Duration:        song.Duration,
		ReplayGain:      song.ReplayGain,
		Peak:            song.Peak,
		AlbumReplayGain: song.AlbumReplayGain,
		AlbumPeak:       song.AlbumPeak,
		ReleaseDate:     song.ReleaseDate,
		TrackNumber:     song.TrackNumber,
		VolumeNumber:    song.VolumeNumber,
		AudioQuality: AudioQuality{
			Name:  song.AudioQuality.Name,
			Color: song.AudioQuality.Color,
		},
		Explicit:   song.Explicit,
		Popularity: song.Popularity,
		Isrc:       song.Isrc,
		Artists:    artists,
		Album: SongAlbum{
			Id:       song.Album.Id,
			Title:    song.Album.Title,
			CoverUrl: song.Album.CoverUrl,
		},
	}
}

type AlbumSongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AlbumSong struct {
	Provider     string            `json:"provider"`
	Downloaded   bool              `json:"downloaded"`
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	TrackNumber  uint              `json:"trackNumber"`
	VolumeNumber uint              `json:"volumeNumber"`
	AudioQuality AudioQuality      `json:"audioQuality"`
	Explicit     bool              `json:"explicit"`
	Isrc         string            `json:"isrc"`
	Artists      []AlbumSongArtist `json:"artists"`
}

type AlbumArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Album struct {
	Provider      string        `json:"provider"`
	Downloaded    bool          `json:"downloaded"`
	Id            string        `json:"id"`
	Title         string        `json:"title"`
	Duration      uint          `json:"duration"`
	ReleaseDate   time.Time     `json:"releaseDate"`
	NumberTracks  uint          `json:"numberTracks"`
	NumberVolumes uint          `json:"numberVolumes"`
	CoverUrl      string        `json:"coverUrl"`
	AudioQuality  AudioQuality  `json:"audioQuality"`
	Explicit      bool          `json:"explicit"`
	Artists       []AlbumArtist `json:"artists"`
	Songs         []AlbumSong   `json:"songs"`
}

func AlbumToResponse(album model.EnrichedAlbum) Album {
	songs := []AlbumSong{}
	for _, s := range album.Songs {
		songArtists := []AlbumSongArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, AlbumSongArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		songs = append(songs, AlbumSong{
			Provider:     s.Provider,
			Downloaded:   s.Downloaded,
			Id:           s.Id,
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

	artists := []AlbumArtist{}
	for _, a := range album.Artists {
		artists = append(artists, AlbumArtist{
			Id:   a.Id,
			Name: a.Name,
		})
	}

	return Album{
		Provider:      album.Provider,
		Downloaded:    album.Downloaded,
		Id:            album.Id,
		Title:         album.Title,
		Duration:      album.Duration,
		ReleaseDate:   album.ReleaseDate,
		NumberTracks:  album.NumberTracks,
		NumberVolumes: album.NumberVolumes,
		CoverUrl:      album.CoverUrl,
		AudioQuality: AudioQuality{
			Name:  album.AudioQuality.Name,
			Color: album.AudioQuality.Color,
		},
		Explicit: album.Explicit,
		Artists:  artists,
		Songs:    songs,
	}
}

type ArtistAlbumArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ArtistAlbum struct {
	Downloaded   bool                `json:"downloaded"`
	Id           string              `json:"id"`
	Title        string              `json:"title"`
	Duration     uint                `json:"duration"`
	ReleaseDate  time.Time           `json:"releaseDate"`
	CoverUrl     string              `json:"coverUrl"`
	AudioQuality AudioQuality        `json:"audioQuality"`
	Explicit     bool                `json:"explicit"`
	Artists      []ArtistAlbumArtist `json:"artists"`
}

type Artist struct {
	Provider   string        `json:"provider"`
	Followed   string        `json:"followed"`
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	PictureUrl string        `json:"pictureUrl"`
	Albums     []ArtistAlbum `json:"albums"`
	Ep         []ArtistAlbum `json:"ep"`
	Singles    []ArtistAlbum `json:"singles"`
}

func ArtistToResponse(artist model.EnrichedArtist) Artist {
	albums := []ArtistAlbum{}
	for _, album := range artist.Albums {
		albumArtists := []ArtistAlbumArtist{}
		for _, a := range album.Artists {
			albumArtists = append(albumArtists, ArtistAlbumArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		albums = append(albums, ArtistAlbum{
			Downloaded:  album.Downloaded,
			Id:          album.Id,
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
	for _, ep := range artist.Ep {
		epArtists := []ArtistAlbumArtist{}
		for _, a := range ep.Artists {
			epArtists = append(epArtists, ArtistAlbumArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		eps = append(eps, ArtistAlbum{
			Downloaded:  ep.Downloaded,
			Id:          ep.Id,
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
	for _, single := range artist.Singles {
		singleArtists := []ArtistAlbumArtist{}
		for _, a := range single.Artists {
			singleArtists = append(singleArtists, ArtistAlbumArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		singles = append(singles, ArtistAlbum{
			Downloaded:  single.Downloaded,
			Id:          single.Id,
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

	return Artist{
		Provider:   artist.Provider,
		Followed:   artist.Followed.String(),
		Id:         artist.Id,
		Name:       artist.Name,
		PictureUrl: artist.PictureUrl,
		Albums:     albums,
		Ep:         eps,
		Singles:    singles,
	}
}

type PlaylistSongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PlaylistSong struct {
	Downloaded   bool                 `json:"downloaded"`
	Id           string               `json:"id"`
	Title        string               `json:"title"`
	Duration     uint                 `json:"duration"`
	AudioQuality AudioQuality         `json:"audioQuality"`
	Explicit     bool                 `json:"explicit"`
	Isrc         string               `json:"isrc"`
	Artists      []PlaylistSongArtist `json:"artists"`
}

type Playlist struct {
	Provider       string         `json:"provider"`
	Downloaded     bool           `json:"downloaded"`
	Id             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Duration       uint           `json:"duration"`
	LastUpdated    time.Time      `json:"lastUpdated"`
	NumberOfTracks uint           `json:"numberOfTracks"`
	CoverURL       string         `json:"coverUrl"`
	Songs          []PlaylistSong `json:"songs"`
}

func PlaylistToResponse(playlist model.EnrichedPlaylist) Playlist {
	songs := []PlaylistSong{}
	for _, s := range playlist.Songs {
		songArtists := []PlaylistSongArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, PlaylistSongArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		songs = append(songs, PlaylistSong{
			Downloaded: s.Downloaded,
			Id:         s.Id,
			Title:      s.Title,
			Duration:   s.Duration,
			AudioQuality: AudioQuality{
				Name:  s.AudioQuality.Name,
				Color: s.AudioQuality.Color,
			},
			Explicit: s.Explicit,
			Isrc:     s.Isrc,
			Artists:  songArtists,
		})
	}
	return Playlist{
		Provider:       playlist.Provider,
		Downloaded:     playlist.Downloaded,
		Id:             playlist.Id,
		Title:          playlist.Title,
		Description:    playlist.Description,
		Duration:       playlist.Duration,
		LastUpdated:    playlist.LastUpdated,
		NumberOfTracks: playlist.NumberOfTracks,
		CoverURL:       playlist.CoverURL,
		Songs:          songs,
	}
}

type SearchItem struct {
	Provider string `json:"provider"`
	Type     string `json:"type"`
	Id       string `json:"id"`
}

type SearchSongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SearchSong struct {
	Downloaded   bool               `json:"downloaded"`
	Id           string             `json:"id"`
	Title        string             `json:"title"`
	Duration     uint               `json:"duration"`
	AudioQuality AudioQuality       `json:"audioQuality"`
	Popularity   uint               `json:"popularity"`
	Explicit     bool               `json:"explicit"`
	Isrc         string             `json:"isrc"`
	Artists      []SearchSongArtist `json:"artists"`
	Album        SearchSongAlbum    `json:"album"`
}

type SearchSongAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type SearchAlbumArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SearchAlbum struct {
	Downloaded   bool                `json:"downloaded"`
	Id           string              `json:"id"`
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
	Id         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}

type SearchPlaylist struct {
	Downloaded bool   `json:"downloaded"`
	Id         string `json:"id"`
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
				Id:   artist.Id,
				Name: artist.Name,
			})
		}

		songs = append(songs, SearchSong{
			Downloaded:   song.Downloaded,
			Id:           song.Id,
			Title:        song.Title,
			Duration:     song.Duration,
			AudioQuality: AudioQuality{Name: song.AudioQuality.Name, Color: song.AudioQuality.Color},
			Popularity:   song.Popularity,
			Explicit:     song.Explicit,
			Isrc:         song.Isrc,
			Artists:      songArtists,
			Album: SearchSongAlbum{
				Id:       song.Album.Id,
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
				Id:   artist.Id,
				Name: artist.Name,
			})
		}

		albums = append(albums, SearchAlbum{
			Downloaded:   album.Downloaded,
			Id:           album.Id,
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
			Id:         artist.Id,
			Name:       artist.Name,
			PictureUrl: artist.PictureUrl,
			Popularity: artist.Popularity,
		})
	}

	playlists := []SearchPlaylist{}
	for _, playlist := range search.Playlists {
		playlists = append(playlists, SearchPlaylist{
			Downloaded: playlist.Downloaded,
			Id:         playlist.Id,
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
				Provider: searchResult.ItemFound.Provider,
				Type:     string(searchResult.ItemFound.Type),
				Id:       searchResult.ItemFound.Id,
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
