package hifi

import "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/dto"

type statusResponse struct {
	Version string `json:"version"`
	Repo    string `json:"Repo"`
}

// -- Song -- //

type songResponse struct {
	Version string       `json:"version"`
	Data    dto.SongItem `json:"data"`
}

// -- Album -- //

type albumResponse struct {
	Version string             `json:"version"`
	Data    albumItemWithSongs `json:"data"`
}

type albumItemWithSongs struct {
	dto.AlbumItem
	Items []dto.SongItemTyped `json:"items"`
}

// -- Artist -- //

type artistResponse struct {
	Version string         `json:"version"`
	Artist  dto.ArtistItem `json:"artist"`
}

type artistAlbumsResponse struct {
	Version string `json:"version"`
	Albums  struct {
		Items []dto.AlbumItem `json:"items"`
	} `json:"albums"`
	Tracks []dto.SongItem `json:"tracks"`
}

// -- Playlist -- //

type playlistResponse struct {
	Version  string           `json:"version"`
	Playlist dto.PlaylistInfo `json:"playlist"`
	Items    []struct {
		dto.SongItemTyped
		DateAdded string `json:"dateAdded"`
		Index     int    `json:"index"`
		ItemUUID  string `json:"itemUuid"`
	} `json:"items"`
}

// -- Search -- //

type searchSongResponse struct {
	Version string `json:"version"`
	Data    struct {
		Limit              int            `json:"limit"`
		Offset             int            `json:"offset"`
		TotalNumberOfItems int            `json:"totalNumberOfItems"`
		Songs              []dto.SongItem `json:"items"`
	} `json:"data"`
}

type searchAlbumResponse struct {
	Version string `json:"version"`
	Data    struct {
		Albums struct {
			Limit              int             `json:"limit"`
			Offset             int             `json:"offset"`
			TotalNumberOfItems int             `json:"totalNumberOfItems"`
			Albums             []dto.AlbumItem `json:"items"`
		} `json:"albums"`
	} `json:"data"`
}

type searchArtistResponse struct {
	Version string `json:"version"`
	Data    struct {
		Artists struct {
			Limit              int              `json:"limit"`
			Offset             int              `json:"offset"`
			TotalNumberOfItems int              `json:"totalNumberOfItems"`
			Artists            []dto.ArtistItem `json:"items"`
		} `json:"artists"`
	} `json:"data"`
}

type searchPlaylistResponse struct {
	Version string `json:"version"`
	Data    struct {
		Playlists struct {
			Limit              int                `json:"limit"`
			Offset             int                `json:"offset"`
			TotalNumberOfItems int                `json:"totalNumberOfItems"`
			Playlists          []dto.PlaylistInfo `json:"items"`
		} `json:"playlists"`
	} `json:"data"`
}

// -- Download -- //

type downloadResponse struct {
	Version string           `json:"version"`
	Data    dto.DownloadItem `json:"data"`
}
