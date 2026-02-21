package hifi

import "github.com/DimitriLaPoudre/MusicShack/server/internal/models"

var (
	LOW = models.Quality{
		Name:  "LOW",
		Color: "#ff0000",
	}
	HIGH = models.Quality{
		Name:  "HIGH",
		Color: "#ff7f00",
	}
	LOSSLESS = models.Quality{
		Name:  "LOSSLESS",
		Color: "#409940",
	}
	HIRES = models.Quality{
		Name:  "HIRES",
		Color: "#00ff00",
	}
)

type status struct {
	Version string `json:"version"`
	Repo    string `json:"Repo"`
}

type artistInfo struct {
	Version string
	Artist  artistItem
	Cover   struct {
		// 750 string
		Id   uint
		Name string
	}
}

type artistAlbums struct {
	Version string
	Albums  struct {
		Items []albumItem
	}
	Tracks []songItem
}

type artistItem struct {
	Id                 uint
	Name               string
	ArtistTypes        []string
	Url                string
	PictureUrl         string `json:"picture"`
	PictureUrlFallback string `json:"selectedAlbumCoverFallback"`
	Popularity         uint
	ArtistRoles        []struct {
		CategoryId int
		Category   string
	}
	Mixes struct {
		Artist_mix string
	}
	Spotlighted bool
}

type albumData struct {
	Version string
	Data    struct {
		Id                     uint
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
			Id   uint
			Name string
			Type string
		}
		Items []struct {
			Item songItem
			Type string
		}
	}
}

type albumItem struct {
	Id                     uint
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
		Id   uint
		Name string
		Type string
	}
}

type albumItemComparaison struct {
	Title       string
	ReleaseDate string
	TrackNumber uint
}

// type albumItemComparaisonExtension struct {
// 	Title       string
// 	ReleaseDate string
// }

type artistMinimalData struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type albumMinimalData struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	CoverUrl     string `json:"cover"`
	VibrantColor string `json:"vibrantColor"`
}

type mixesMinimalData struct {
	Track_mix string `json:"TRACK_MIX"`
}

type playlistInfoData struct {
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

type playlistData struct {
	Version  string           `json:"version"`
	Playlist playlistInfoData `json:"playlist"`
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
			Mixes       mixesMinimalData    `json:"mixes"`
			DateAdded   string              `json:"dateAdded"`
			Index       uint                `json:"index"`
			ItemUUID    string              `json:"itemUuid"`
		} `json:"item"`
	} `json:"items"`
}

type songData struct {
	Version string
	Data    songItem
}

type songItem struct {
	Id                     uint
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
		Id   uint
		Name string
		Type string
	}
	Artists []struct {
		Id   uint
		Name string
		Type string
	}
	Album struct {
		Id           uint
		Title        string
		CoverUrl     string `json:"cover"`
		VibrantColor string
	}
	Mixes struct {
		Track_mix string
	}
	CoverUrl string
}

type searchSongData struct {
	Version string
	Data    struct {
		Limit              uint
		Offset             uint
		TotalNumberOfItems uint
		Songs              []songItem `json:"items"`
	}
}

type searchAlbumData struct {
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

type searchArtistData struct {
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

type searchPlaylistData struct {
	Version string `json:"version"`
	Data    struct {
		Playlists struct {
			Limit              uint               `json:"limit"`
			Offset             uint               `json:"offset"`
			TotalNumberOfItems uint               `json:"totalNumberOfItems"`
			Playlists          []playlistInfoData `json:"items"`
		} `json:"playlists"`
	} `json:"data"`
}

type downloadData struct {
	Version string
	Data    downloadItem
}

type downloadItem struct {
	TrackId            uint
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
