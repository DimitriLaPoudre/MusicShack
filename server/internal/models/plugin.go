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
		Id    string
		Title string
	}
	BitDepth    uint
	SampleRate  uint
	DownloadUrl string
}
