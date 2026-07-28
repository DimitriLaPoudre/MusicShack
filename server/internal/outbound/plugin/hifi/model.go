package hifi

import "github.com/DimitriLaPoudre/MusicShack/internal/model"

const (
	AudioQualityLOW      string = "LOW"
	AudioQualityHIGH     string = "HIGH"
	AudioQualityLOSSLESS string = "LOSSLESS"
	AudioQualityHIRES    string = "HI_RES_LOSSLESS"
)

var (
	LOW = model.Quality{
		Name:  "LOW",
		Color: "#ff0000",
	}
	HIGH = model.Quality{
		Name:  "HIGH",
		Color: "#ff7f00",
	}
	LOSSLESS = model.Quality{
		Name:  "LOSSLESS",
		Color: "#409940",
	}
	HIRES = model.Quality{
		Name:  "HIRES",
		Color: "#00ff00",
	}
)

const StreamStartDateLayout = "2006-01-02T15:04:05.000Z0700"
const ReleaseDateLayout = "2006-01-02"

type statusResponse struct {
	Version string `json:"version"`
	Repo    string `json:"Repo"`
}

// -- Song -- //

type songResponse struct {
	Version string   `json:"version"`
	Data    songItem `json:"data"`
}

type songItem struct {
	ID                     uint64  `json:"id"`
	Title                  string  `json:"title"`
	Duration               uint    `json:"duration"`
	ReplayGain             float64 `json:"replayGain"`
	Peak                   float64 `json:"peak"`
	AllowStreaming         bool    `json:"allowStreaming"`
	StreamReady            bool    `json:"streamReady"`
	PayToStream            bool    `json:"payToStream"`
	AdSupportedStreamReady bool    `json:"adSupportedStreamReady"`
	DjReady                bool    `json:"djReady"`
	StemReady              bool    `json:"stemReady"`
	StreamStartDate        string  `json:"streamStartDate"`
	PremiumStreamingOnly   bool    `json:"premiumStreamingOnly"`
	TrackNumber            uint    `json:"trackNumber"`
	VolumeNumber           uint    `json:"volumeNumber"`

	// Version *string `json:"version"`

	Popularity uint   `json:"popularity"`
	Copyright  string `json:"copyright"`

	// BPM      *uint   `json:"bpm"`
	// Key      *string `json:"key"`
	// KeyScale *string `json:"keyScale"`

	URL      string `json:"url"`
	ISRC     string `json:"isrc"`
	Editable bool   `json:"editable"`
	Explicit bool   `json:"explicit"`

	AudioQuality string   `json:"audioQuality"`
	AudioModes   []string `json:"audioModes"`

	MediaMetadata struct {
		Tags []string `json:"tags"`
	} `json:"mediaMetadata"`

	Upload      bool    `json:"upload"`
	AccessType  *string `json:"accessType"`
	Spotlighted bool    `json:"spotlighted"`
	AI          bool    `json:"ai"`

	Artist  miniArtistItem   `json:"artist"`
	Artists []miniArtistItem `json:"artists"`

	Album struct {
		ID           uint64 `json:"id"`
		Title        string `json:"title"`
		Cover        string `json:"cover"`
		VibrantColor string `json:"vibrantColor"`

		// VideoCover *string `json:"videoCover"`
	} `json:"album"`

	Mixes struct {
		TrackMix string `json:"TRACK_MIX"`
	} `json:"mixes"`
}

type SongItemTyped struct {
	Item songItem `json:"item"`
	Type string   `json:"type"`
}

// -- Album -- //

type albumResponse struct {
	Version string             `json:"version"`
	Data    albumItemWithSongs `json:"data"`
}

type albumItemWithSongs struct {
	albumItem
	Items []SongItemTyped `json:"items"`
}

type albumItem struct {
	ID                     uint64 `json:"id"`
	Title                  string `json:"title"`
	Duration               uint   `json:"duration"`
	StreamReady            bool   `json:"streamReady"`
	PayToStream            bool   `json:"payToStream"`
	AdSupportedStreamReady bool   `json:"adSupportedStreamReady"`
	DjReady                bool   `json:"djReady"`
	StemReady              bool   `json:"stemReady"`
	StreamStartDate        string `json:"streamStartDate"`
	AllowStreaming         bool   `json:"allowStreaming"`
	PremiumStreamingOnly   bool   `json:"premiumStreamingOnly"`

	NumberOfTracks  uint `json:"numberOfTracks"`
	NumberOfVideos  uint `json:"numberOfVideos"`
	NumberOfVolumes uint `json:"numberOfVolumes"`

	ReleaseDate string `json:"releaseDate"`
	Copyright   string `json:"copyright"`
	Type        string `json:"type"`

	// Version *string `json:"version"`

	URL          string `json:"url"`
	Cover        string `json:"cover"`
	VibrantColor string `json:"vibrantColor"`

	// VideoCover *string `json:"videoCover"`

	Explicit     bool     `json:"explicit"`
	UPC          string   `json:"upc"`
	Popularity   uint     `json:"popularity"`
	AudioQuality string   `json:"audioQuality"`
	AudioModes   []string `json:"audioModes"`

	MediaMetadata struct {
		Tags []string `json:"tags"`
	} `json:"mediaMetadata"`

	Upload bool `json:"upload"`
	AI     bool `json:"ai"`

	Artist miniArtistItem `json:"artist"` // struct in /album but not /search/?al

	Artists []miniArtistItem `json:"artists"`
}

// -- Artist -- //

type artistResponse struct {
	Version string `json:"version"`

	Artist artistItem `json:"artist"`

	Cover coverItem `json:"cover"`
}

type artistItem struct {
	ID                         uint64   `json:"id"`
	Name                       string   `json:"name"`
	ArtistTypes                []string `json:"artistTypes"`
	URL                        string   `json:"url"`
	Picture                    string   `json:"picture"`
	SelectedAlbumCoverFallback *string  `json:"selectedAlbumCoverFallback"`
	Popularity                 uint     `json:"popularity"`

	ArtistRoles []struct {
		CategoryID int    `json:"categoryId"`
		Category   string `json:"category"`
	} `json:"artistRoles"`

	Mixes struct {
		ArtistMix string `json:"ARTIST_MIX"`
	} `json:"mixes"`

	// Handle *string `json:"handle"`
	// UserID *uint64 `json:"userId"`

	Spotlighted bool `json:"spotlighted"`
}

type miniArtistItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`

	// Handle *string `json:"handle"`

	Type string `json:"type"`

	Picture *string `json:"picture"`
}

type artistAlbumsResponse struct {
	Version string `json:"version"`
	Albums  struct {
		Items []albumItem `json:"items"`
	} `json:"albums"`
	Tracks []songItem `json:"tracks"`
}

// -- Playlist -- //

type playlistResponse struct {
	Version  string       `json:"version"`
	Playlist playlistInfo `json:"playlist"`
	Items    []struct {
		SongItemTyped
		DateAdded string `json:"dateAdded"`
		Index     uint   `json:"index"`
		ItemUUID  string `json:"itemUuid"`
	} `json:"items"`
}

type playlistInfo struct {
	UUID           string `json:"uuid"`
	Title          string `json:"title"`
	NumberOfTracks uint   `json:"numberOfTracks"`
	NumberOfVideos uint   `json:"numberOfVideos"`
	Creator        struct {
		ID uint `json:"id"`
	} `json:"creator"`
	Description     string           `json:"description"`
	Duration        uint             `json:"duration"`
	LastUpdated     string           `json:"lastUpdated"`
	Created         string           `json:"created"`
	Type            string           `json:"type"`
	PublicPlaylist  bool             `json:"publicPlaylist"`
	URL             string           `json:"url"`
	Image           string           `json:"image"`
	Popularity      uint             `json:"popularity"`
	SquareImage     string           `json:"squareImage"`
	PromotedArtists []miniArtistItem `json:"promotedArtists"`
	LastItemAddedAt string           `json:"lastItemAddedAt"`
}

// -- Search -- //

type searchSongResponse struct {
	Version string `json:"version"`
	Data    struct {
		Limit              uint       `json:"limit"`
		Offset             uint       `json:"offset"`
		TotalNumberOfItems uint       `json:"totalNumberOfItems"`
		Songs              []songItem `json:"items"`
	} `json:"data"`
}

type searchAlbumResponse struct {
	Version string `json:"version"`
	Data    struct {
		Albums struct {
			Limit              uint        `json:"limit"`
			Offset             uint        `json:"offset"`
			TotalNumberOfItems uint        `json:"totalNumberOfItems"`
			Albums             []albumItem `json:"items"`
		} `json:"albums"`
	} `json:"data"`
}

type searchArtistResponse struct {
	Version string `json:"version"`
	Data    struct {
		Artists struct {
			Limit              uint         `json:"limit"`
			Offset             uint         `json:"offset"`
			TotalNumberOfItems uint         `json:"totalNumberOfItems"`
			Artists            []artistItem `json:"items"`
		} `json:"artists"`
	} `json:"data"`
}

type searchPlaylistResponse struct {
	Version string `json:"version"`
	Data    struct {
		Playlists struct {
			Limit              uint           `json:"limit"`
			Offset             uint           `json:"offset"`
			TotalNumberOfItems uint           `json:"totalNumberOfItems"`
			Playlists          []playlistInfo `json:"items"`
		} `json:"playlists"`
	} `json:"data"`
}

// -- Cover -- //

type coverItem struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Size string `json:"750"`
}

// -- Download -- //

type downloadResponse struct {
	Version string       `json:"version"`
	Data    downloadItem `json:"data"`
}

type downloadItem struct {
	TrackID            uint
	AssetPresentation  string
	AudioMode          string
	AudioQuality       string
	ManifestMimeType   string
	ManifestHash       string
	Manifest           string
	AlbumReplayGain    float64
	AlbumPeakAmplitude float64
	TrackReplayGain    float64
	TrackPeakAmplitude float64
	BitDepth           uint
	SampleRate         uint
}

type manifestTidal struct {
	MimeType       string
	Codecs         string
	EncryptionType string
	Urls           []string
}

type manifestMPD struct {
	Periods []Period `xml:"Period"`
}

type Period struct {
	AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
	Representations []Representation `xml:"Representation"`
}

type Representation struct {
	SegmentTemplate SegmentTemplate `xml:"SegmentTemplate"`
}

type SegmentTemplate struct {
	Initialization string          `xml:"initialization,attr"`
	Media          string          `xml:"media,attr"`
	StartNumber    int             `xml:"startNumber,attr"`
	Timeline       SegmentTimeline `xml:"SegmentTimeline"`
}

type SegmentTimeline struct {
	Segments []Segment `xml:"S"`
}

type Segment struct {
	D int `xml:"d,attr"`
	R int `xml:"r,attr"`
}
