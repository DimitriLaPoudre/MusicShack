package response

import "github.com/Ascension-EIP/Ascension/apps/server/internal/model"

type Quality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type MiniSong struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	AudioQuality Quality      `json:"audioQuality"`
	Explicit     bool         `json:"explicit"`
	Isrc         string       `json:"isrc"`
	Artists      []MiniArtist `json:"artists"`
}

type Song struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	TrackNumber  uint         `json:"trackNumber"`
	VolumeNumber uint         `json:"volumeNumber"`
	AudioQuality Quality      `json:"audioQuality"`
	Explicit     bool         `json:"explicit"`
	Isrc         string       `json:"isrc"`
	Artists      []MiniArtist `json:"artists"`
}

type FullSong struct {
	Provider        string       `json:"provider"`
	Downloaded      bool         `json:"downloaded"`
	Id              string       `json:"id"`
	Title           string       `json:"title"`
	Duration        uint         `json:"duration"`
	ReplayGain      float64      `json:"replayGain"`
	Peak            float64      `json:"peak"`
	AlbumReplayGain float64      `json:"albumReplayGain"`
	AlbumPeak       float64      `json:"albumPeak"`
	ReleaseDate     string       `json:"releaseDate"`
	TrackNumber     uint         `json:"trackNumber"`
	VolumeNumber    uint         `json:"volumeNumber"`
	AudioQuality    Quality      `json:"audioQuality"`
	Explicit        bool         `json:"explicit"`
	Popularity      uint         `json:"popularity"`
	Isrc            string       `json:"isrc"`
	Artists         []MiniArtist `json:"artists"`
	Album           MiniAlbum    `json:"album"`
}

func SongToResponse(song model.Song) FullSong {
	artists := []MiniArtist{}
	for _, a := range song.Artists {
		artists = append(artists, MiniArtist{
			Id:   a.Id,
			Name: a.Name,
		})
	}

	return FullSong{
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
		AudioQuality: Quality{
			Name:  song.AudioQuality.Name,
			Color: song.AudioQuality.Color,
		},
		Explicit:   song.Explicit,
		Popularity: song.Popularity,
		Isrc:       song.Isrc,
		Artists:    artists,
		Album: MiniAlbum{
			Id:       song.Album.Id,
			Title:    song.Album.Title,
			CoverUrl: song.Album.CoverUrl,
		},
	}
}

type MiniAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type Album struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	ReleaseDate  string       `json:"releaseDate"`
	CoverUrl     string       `json:"coverUrl"`
	AudioQuality Quality      `json:"audioQuality"`
	Explicit     bool         `json:"explicit"`
	Artists      []MiniArtist `json:"artists"`
}

type FullAlbum struct {
	Provider      string       `json:"provider"`
	Downloaded    bool         `json:"downloaded"`
	Id            string       `json:"id"`
	Title         string       `json:"title"`
	Duration      uint         `json:"duration"`
	ReleaseDate   string       `json:"releaseDate"`
	NumberTracks  uint         `json:"numberTracks"`
	NumberVolumes uint         `json:"numberVolumes"`
	CoverUrl      string       `json:"coverUrl"`
	AudioQuality  Quality      `json:"audioQuality"`
	Explicit      bool         `json:"explicit"`
	Artists       []MiniArtist `json:"artists"`
	Songs         []Song       `json:"songs"`
}

func AlbumToResponse(album model.Album) FullAlbum {
	songs := []Song{}
	for _, s := range album.Songs {
		songArtists := []MiniArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, MiniArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		songs = append(songs, Song{
			Downloaded:   s.Downloaded,
			Id:           s.Id,
			Title:        s.Title,
			Duration:     s.Duration,
			TrackNumber:  s.TrackNumber,
			VolumeNumber: s.VolumeNumber,
			AudioQuality: Quality{
				Name:  s.AudioQuality.Name,
				Color: s.AudioQuality.Color,
			},
			Explicit: s.Explicit,
			Isrc:     s.Isrc,
			Artists:  songArtists,
		})
	}

	artists := []MiniArtist{}
	for _, a := range album.Artists {
		artists = append(artists, MiniArtist{
			Id:   a.Id,
			Name: a.Name,
		})
	}

	return FullAlbum{
		Provider:      album.Provider,
		Downloaded:    album.Downloaded,
		Id:            album.Id,
		Title:         album.Title,
		Duration:      album.Duration,
		ReleaseDate:   album.ReleaseDate,
		NumberTracks:  album.NumberTracks,
		NumberVolumes: album.NumberVolumes,
		CoverUrl:      album.CoverUrl,
		AudioQuality: Quality{
			Name:  album.AudioQuality.Name,
			Color: album.AudioQuality.Color,
		},
		Explicit: album.Explicit,
		Artists:  artists,
		Songs:    songs,
	}
}

type MiniArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FullArtist struct {
	Provider   string  `json:"provider"`
	Followed   uint    `json:"followed"`
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	PictureUrl string  `json:"pictureUrl"`
	Albums     []Album `json:"albums"`
	Ep         []Album `json:"ep"`
	Singles    []Album `json:"singles"`
}

func ArtistToResponse(artist model.Artist) FullArtist {
	albums := []Album{}
	for _, album := range artist.Albums {
		albumArtists := []MiniArtist{}
		for _, a := range album.Artists {
			albumArtists = append(albumArtists, MiniArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		albums = append(albums, Album{
			Downloaded:  album.Downloaded,
			Id:          album.Id,
			Title:       album.Title,
			Duration:    album.Duration,
			ReleaseDate: album.ReleaseDate,
			CoverUrl:    album.CoverUrl,
			AudioQuality: Quality{
				Name:  album.AudioQuality.Name,
				Color: album.AudioQuality.Color,
			},
			Explicit: album.Explicit,
			Artists:  albumArtists,
		})
	}

	eps := []Album{}
	for _, ep := range artist.Ep {
		epArtists := []MiniArtist{}
		for _, a := range ep.Artists {
			epArtists = append(epArtists, MiniArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		eps = append(eps, Album{
			Downloaded:  ep.Downloaded,
			Id:          ep.Id,
			Title:       ep.Title,
			Duration:    ep.Duration,
			ReleaseDate: ep.ReleaseDate,
			CoverUrl:    ep.CoverUrl,
			AudioQuality: Quality{
				Name:  ep.AudioQuality.Name,
				Color: ep.AudioQuality.Color,
			},
			Explicit: ep.Explicit,
			Artists:  epArtists,
		})
	}

	singles := []Album{}
	for _, single := range artist.Singles {
		singleArtists := []MiniArtist{}
		for _, a := range single.Artists {
			singleArtists = append(singleArtists, MiniArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		singles = append(singles, Album{
			Downloaded:  single.Downloaded,
			Id:          single.Id,
			Title:       single.Title,
			Duration:    single.Duration,
			ReleaseDate: single.ReleaseDate,
			CoverUrl:    single.CoverUrl,
			AudioQuality: Quality{
				Name:  single.AudioQuality.Name,
				Color: single.AudioQuality.Color,
			},
			Explicit: single.Explicit,
			Artists:  singleArtists,
		})
	}

	return FullArtist{
		Provider:   artist.Provider,
		Followed:   artist.Followed,
		Id:         artist.Id,
		Name:       artist.Name,
		PictureUrl: artist.PictureUrl,
		Albums:     albums,
		Ep:         eps,
		Singles:    singles,
	}
}

type FullPlaylist struct {
	Provider       string     `json:"provider"`
	Downloaded     bool       `json:"downloaded"`
	Id             string     `json:"id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Duration       uint       `json:"duration"`
	LastUpdated    string     `json:"lastUpdated"`
	NumberOfTracks uint       `json:"numberOfTracks"`
	CoverURL       string     `json:"coverUrl"`
	Songs          []MiniSong `json:"songs"`
}

func PlaylistToResponse(playlist model.Playlist) FullPlaylist {
	songs := []MiniSong{}
	for _, s := range playlist.Songs {
		songArtists := []MiniArtist{}
		for _, a := range s.Artists {
			songArtists = append(songArtists, MiniArtist{
				Id:   a.Id,
				Name: a.Name,
			})
		}

		songs = append(songs, MiniSong{
			Downloaded: s.Downloaded,
			Id:         s.Id,
			Title:      s.Title,
			Duration:   s.Duration,
			AudioQuality: Quality{
				Name:  s.AudioQuality.Name,
				Color: s.AudioQuality.Color,
			},
			Explicit: s.Explicit,
			Isrc:     s.Isrc,
			Artists:  songArtists,
		})
	}
	return FullPlaylist{
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
