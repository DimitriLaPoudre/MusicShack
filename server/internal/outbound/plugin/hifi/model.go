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

type songItem struct {
	ID                     uint
	Title                  string
	Duration               uint
	ReplayGain             float64
	Peak                   float64
	AllowStreaming         bool
	StreamReady            bool
	PayToStream            bool
	AdSupportedStreamReady bool
	DjReady                bool
	StemReady              bool
	ReleaseDate            string `json:"StreamStartDate"`
	PremiumStreamingOnly   bool
	TrackNumber            uint
	VolumeNumber           uint
	Popularity             uint
	Copyright              string
	Bpm                    uint
	Key                    string
	KeyScale               string
	Url                    string
	Isrc                   string
	Editable               bool
	Explicit               bool
	AudioQuality           string
	AudioModes             []string
	MediaMetadata          struct {
		Tags []string
	}
	Upload      bool
	AccessType  string
	Spotlighted bool
	Artist      struct {
		ID   uint
		Name string
		Type string
	}
	Artists []struct {
		ID   uint
		Name string
		Type string
	}
	Album struct {
		ID           uint
		Title        string
		CoverUrl     string `json:"cover"`
		VibrantColor string
	}
	Mixes struct {
		Track_mix string
	}
	CoverUrl string
}

type songResponse struct {
	Version string
	Data    songItem
}

type albumMinimalData struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	CoverUrl     string `json:"cover"`
	VibrantColor string `json:"vibrantColor"`
}

type albumItem struct {
	ID                     uint
	Title                  string
	Duration               uint
	StreamReady            bool
	PayToStream            bool
	AdSupportedStreamReady bool
	DjReady                bool
	StemReady              bool
	StreamStartDate        string
	AllowStreaming         bool
	PremiumStreamingOnly   bool
	NumberOfTracks         uint
	NumberOfVideos         uint
	NumberOfVolumes        uint
	ReleaseDate            string
	Copyright              string
	Type                   string
	Url                    string
	CoverUrl               string `json:"cover"`
	VibrantColor           string
	Explicit               bool
	Upc                    string
	Popularity             uint
	AudioQuality           string
	AudioModes             []string
	MediaMetadata          struct {
		Tags []string
	}
	Upload  bool
	Artists []struct {
		ID   uint
		Name string
		Type string
	}
}

type albumResponse struct {
	Version string
	Data    struct {
		ID                     uint
		Title                  string
		Duration               uint
		StreamReady            bool
		PayToStream            bool
		AdSupportedStreamReady bool
		DjReady                bool
		StemReady              bool
		StreamStartDate        string
		AllowStreaming         bool
		PremiumStreamingOnly   bool
		NumberOfTracks         uint
		NumberOfVideos         uint
		NumberOfVolumes        uint
		ReleaseDate            string
		Copyright              string
		Type                   string
		Url                    string
		CoverUrl               string `json:"cover"`
		VibrantColor           string
		Explicit               bool
		Upc                    string
		Popularity             uint
		AudioQuality           string
		AudioModes             []string
		MediaMetadata          struct {
			Tags []string
		}
		Upload  bool
		Artists []struct {
			ID   uint
			Name string
			Type string
		}
		Items []struct {
			Item songItem
			Type string
		}
	}
}

type artistMinimalData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type artistItem struct {
	ID                 uint
	Name               string
	ArtistTypes        []string
	Url                string
	PictureUrl         string `json:"picture"`
	PictureUrlFallback string `json:"selectedAlbumCoverFallback"`
	Popularity         uint
	ArtistRoles        []struct {
		CategoryID int
		Category   string
	}
	Mixes struct {
		Artist_mix string
	}
	Spotlighted bool
}

type artistResponse struct {
	Version string
	Artist  artistItem
	Cover   struct {
		// 750 string
		ID   uint
		Name string
	}
}

type artistAlbumsResponse struct {
	Version string
	Albums  struct {
		Items []albumItem
	}
	Tracks []songItem
}

type playlistInfo struct {
	UUID           string `json:"uuid"`
	Title          string `json:"title"`
	NumberOfTracks uint   `json:"numberOfTracks"`
	NumberOfVideos uint   `json:"numberOfVideos"`
	Creator        struct {
		ID uint `json:"id"`
	} `json:"creator"`
	Description     string              `json:"description"`
	Duration        uint                `json:"duration"`
	LastUpdated     string              `json:"lastUpdated"`
	Created         string              `json:"created"`
	Type            string              `json:"type"`
	PublicPlaylist  bool                `json:"publicPlaylist"`
	URL             string              `json:"url"`
	Image           string              `json:"image"`
	Popularity      uint                `json:"popularity"`
	SquareImage     string              `json:"squareImage"`
	PromotedArtists []artistMinimalData `json:"promotedArtists"`
	LastItemAddedAt string              `json:"lastItemAddedAt"`
}

type playlistResponse struct {
	Version  string       `json:"version"`
	Playlist playlistInfo `json:"playlist"`
	Items    []struct {
		Type string `json:"type"`
		Item struct {
			ID                     uint     `json:"id"`
			Title                  string   `json:"title"`
			Duration               uint     `json:"duration"`
			ReplayGain             float64  `json:"replayGain"`
			Peak                   float64  `json:"peak"`
			AllowStreaming         bool     `json:"allowStreaming"`
			StreamReady            bool     `json:"streamReady"`
			PayToStream            bool     `json:"payToStream"`
			AdSupportedStreamReady bool     `json:"adSupportedStreamReady"`
			DjReady                bool     `json:"djReady"`
			StemReady              bool     `json:"stemReady"`
			ReleaseDate            string   `json:"streamStartDate"`
			PremiumStreamingOnly   bool     `json:"premiumStreamingOnly"`
			TrackNumber            uint     `json:"trackNumber"`
			VolumeNumber           uint     `json:"volumeNumber"`
			Popularity             uint     `json:"popularity"`
			Copyright              string   `json:"copyright"`
			Bpm                    uint     `json:"bpm"`
			Key                    string   `json:"key"`
			KeyScale               string   `json:"keyScale"`
			URL                    string   `json:"url"`
			ISRC                   string   `json:"isrc"`
			Editable               bool     `json:"editable"`
			Explicit               bool     `json:"explicit"`
			AudioQuality           string   `json:"audioQuality"`
			AudioModes             []string `json:"audioModes"`
			MediaMetadata          struct {
				Tags []string `json:"tags"`
			} `json:"mediaMetadata"`
			Upload      bool                `json:"upload"`
			AccessType  string              `json:"accessType"`
			Spotlighted bool                `json:"spotlighted"`
			Artist      artistMinimalData   `json:"artist"`
			Artists     []artistMinimalData `json:"artists"`
			Album       albumMinimalData    `json:"album"`
			Mixes       struct {
				Track_mix string `json:"TRACK_MIX"`
			} `json:"mixes"`
			DateAdded string `json:"dateAdded"`
			Index     uint   `json:"index"`
			ItemUUID  string `json:"itemUuid"`
		} `json:"item"`
	} `json:"items"`
}

type searchSongResponse struct {
	Version string
	Data    struct {
		Limit              uint
		Offset             uint
		TotalNumberOfItems uint
		Songs              []songItem `json:"items"`
	}
}

type searchAlbumResponse struct {
	Version string
	Data    struct {
		Albums struct {
			Limit              uint
			Offset             uint
			TotalNumberOfItems uint
			Albums             []albumItem `json:"items"`
		}
	}
}

type searchArtistResponse struct {
	Version string
	Data    struct {
		Artists struct {
			Limit              uint
			Offset             uint
			TotalNumberOfItems uint
			Artists            []artistItem `json:"items"`
		}
	}
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

type downloadResponse struct {
	Version string
	Data    downloadItem
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
