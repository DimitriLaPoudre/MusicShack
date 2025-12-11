package hifiv2

type artistData struct {
	Artist struct {
		Id                 uint
		Name               string
		PictureUrl         string
		PictureUrlFallback string `json:"selectedAlbumCoverFallback"`
	}
}

type artistAlbums struct {
	Albums struct {
		Id    string
		Title string
		Rows  []struct {
			Modules []struct {
				PagedList struct {
					Limit               uint
					Offset              uint
					TotalNumbersOfItems uint
					Items               []struct {
						Id       uint
						Title    string
						CoverUrl string `json:"cover"`
						Artists  []struct {
							Id   uint
							Name string
						}
					}
				}
			}
		}
	}
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
}
