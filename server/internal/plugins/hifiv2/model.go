package hifiv2

type artistData struct {
	Artist struct {
		Id                 uint   `mapstructure:"id"`
		Name               string `mapstructure:"name"`
		PictureUrl         string `mapstructure:"picture"`
		PictureUrlFallback string `mapstructure:"selectedAlbumCoverFallback"`
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
						CoverUrl string `mapstructure:"cover"`
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
	Id            uint   `mapstructure:"id"`
	Title         string `mapstructure:"title"`
	Duration      uint   `mapstructure:"duration"`
	ReleaseDate   string `mapstructure:"releaseDate"`
	NumberTracks  uint   `mapstructure:"numberOfTracks"`
	NumberVolumes uint   `mapstructure:"numberOfVolumes"`
	Type          string `mapstructure:"type"`
	CoverUrl      string `mapstructure:"cover"`
	AudioQuality  string `mapstructure:"audioQuality"`
	Artist        struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artist"`
	Artists []struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artists"`
	Limit       uint `mapstructure:"limit"`
	Offset      uint `mapstructure:"offset"`
	NumberSongs uint `mapstructure:"totalNumberOfItems"`
	DirtySongs  []struct {
		SongData struct {
			Id           uint   `mapstructure:"id"`
			Title        string `mapstructure:"title"`
			Duration     uint   `mapstructure:"duration"`
			TrackNumber  uint   `mapstructure:"trackNumber"`
			VolumeNumber uint   `mapstructure:"volumeNumber"`
			Artists      []struct {
				Id   uint   `mapstructure:"id"`
				Name string `mapstructure:"name"`
			} `mapstructure:"artists"`
		} `mapstructure:"item"`
		Type string `mapstructure:"type"`
	} `mapstructure:"items"`
	Songs []struct {
		Id           uint
		Title        string
		Duration     uint
		TrackNumber  uint
		VolumeNumber uint
		Artists      []struct {
			Id   uint
			Name string
		}
	}
}

type songData struct {
	Data struct {
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
}
