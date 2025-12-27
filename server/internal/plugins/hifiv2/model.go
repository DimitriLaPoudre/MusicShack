package hifiv2

type artistData struct {
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
		SelfLink string
		Id       string
		Title    string
		Rows     []struct {
			Modules []struct {
				PagedList struct {
					DataApiPath         string
					Limit               uint
					Offset              uint
					TotalNumbersOfItems uint
					Items               []albumItem
				}
			}
		}
	}
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
		Limit              uint
		Offset             uint
		TotalNumberOfItems uint
		Items              []struct {
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
