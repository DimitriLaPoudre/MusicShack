package models

type SongData struct {
	Id           string
	Title        string
	Duration     uint
	ReleaseDate  string
	TrackNumber  uint
	VolumeNumber uint
	AudioQuality string
	Artist       struct {
		Id   string
		Name string
	}
	Artists []struct {
		Id   string
		Name string
	}
	Album struct {
		Id       string
		Title    string
		CoverUrl string
	}
	BitDepth    uint
	SampleRate  uint
	DownloadUrl string
}

type AlbumData struct {
	Id            string
	Title         string
	Duration      uint
	ReleaseDate   string
	NumberTracks  uint
	NumberVolumes uint
	Type          string
	CoverUrl      string
	AudioQuality  string
	Artist        struct {
		Id   string
		Name string
	}
	Artists []struct {
		Id   string
		Name string
	}
	Limit       uint
	Offset      uint
	NumberSongs uint
	Songs       []struct {
		Id           string
		Title        string
		Duration     uint
		TrackNumber  uint
		VolumeNumber uint
		Artists      []struct {
			Id   string
			Name string
		}
	}
}

type ArtistData struct {
	Id         string
	Name       string
	PictureUrl string
	Albums     []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
	Ep []struct {
		Id       string
		Title    string
		CoverUrl string
	}
	Singles []struct {
		Id       string
		Title    string
		CoverUrl string
	}
}

type SearchData struct {
	Songs []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
	Albums []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
	Artists []struct {
		Id         string
		Name       string
		PictureUrl string
	}
}
