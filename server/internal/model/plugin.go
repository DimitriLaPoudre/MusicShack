package model

type Quality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Song struct {
	Provider        string       `json:"provider"`
	Api             string       `json:"api"`
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
	Artists         []SongArtist `json:"artists"`
	Album           SongAlbum    `json:"album"`
}

type SongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SongAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type Playlist struct {
	Provider       string         `json:"provider"`
	Api            string         `json:"api"`
	Downloaded     bool           `json:"downloaded"`
	Id             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Duration       uint           `json:"duration"`
	LastUpdated    string         `json:"lastUpdated"`
	NumberOfTracks uint           `json:"numberOfTracks"`
	CoverURL       string         `json:"coverUrl"`
	Songs          []PlaylistSong `json:"songs"`
}

type PlaylistSong struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	AudioQuality Quality      `json:"audioQuality"`
	Explicit     bool         `json:"explicit"`
	Isrc         string       `json:"isrc"`
	Artists      []SongArtist `json:"artists"`
}

type Album struct {
	Provider      string        `json:"provider"`
	Api           string        `json:"api"`
	Downloaded    bool          `json:"downloaded"`
	Id            string        `json:"id"`
	Title         string        `json:"title"`
	Duration      uint          `json:"duration"`
	ReleaseDate   string        `json:"releaseDate"`
	NumberTracks  uint          `json:"numberTracks"`
	NumberVolumes uint          `json:"numberVolumes"`
	CoverUrl      string        `json:"coverUrl"`
	AudioQuality  Quality       `json:"audioQuality"`
	Explicit      bool          `json:"explicit"`
	Artists       []AlbumArtist `json:"artists"`
	Songs         []AlbumSong   `json:"songs"`
}

type AlbumArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AlbumSong struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	TrackNumber  uint         `json:"trackNumber"`
	VolumeNumber uint         `json:"volumeNumber"`
	AudioQuality Quality      `json:"audioQuality"`
	Explicit     bool         `json:"explicit"`
	Isrc         string       `json:"isrc"`
	Artists      []SongArtist `json:"artists"`
}

type Artist struct {
	Provider   string        `json:"provider"`
	Api        string        `json:"api"`
	Followed   uint          `json:"followed"`
	Id         string        `json:"id"`
	Name       string        `json:"name"`
	PictureUrl string        `json:"pictureUrl"`
	Albums     []ArtistAlbum `json:"albums"`
	Ep         []ArtistAlbum `json:"ep"`
	Singles    []ArtistAlbum `json:"singles"`
}

type ArtistAlbum struct {
	Downloaded   bool          `json:"downloaded"`
	Id           string        `json:"id"`
	Title        string        `json:"title"`
	Duration     uint          `json:"duration"`
	ReleaseDate  string        `json:"releaseDate"`
	CoverUrl     string        `json:"coverUrl"`
	AudioQuality Quality       `json:"audioQuality"`
	Explicit     bool          `json:"explicit"`
	Artists      []AlbumArtist `json:"artists"`
}

type Search struct {
	Songs     []SearchSong     `json:"songs"`
	Albums    []SearchAlbum    `json:"albums"`
	Artists   []SearchArtist   `json:"artists"`
	Playlists []SearchPlaylist `json:"playlists"`
}

type SearchSong struct {
	Downloaded   bool         `json:"downloaded"`
	Id           string       `json:"id"`
	Title        string       `json:"title"`
	Duration     uint         `json:"duration"`
	AudioQuality Quality      `json:"audioQuality"`
	Popularity   uint         `json:"popularity"`
	Explicit     bool         `json:"explicit"`
	Isrc         string       `json:"isrc"`
	Artists      []SongArtist `json:"artists"`
	Album        SongAlbum    `json:"album"`
}

type SearchAlbum struct {
	Downloaded   bool          `json:"downloaded"`
	Id           string        `json:"id"`
	Title        string        `json:"title"`
	Duration     uint          `json:"duration"`
	CoverUrl     string        `json:"coverUrl"`
	AudioQuality Quality       `json:"audioQuality"`
	Explicit     bool          `json:"explicit"`
	Popularity   uint          `json:"popularity"`
	Artists      []AlbumArtist `json:"artists"`
}

type SearchArtist struct {
	Followed   uint   `json:"followed"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}

type SearchPlaylist struct {
	Downloaded bool   `json:"downloaded"`
	ID         string `json:"id"`
	Title      string `json:"title"`
	Duration   uint   `json:"duration"`
	CoverURL   string `json:"coverUrl"`
	Popularity uint   `json:"popularity"`
}

type Type string

const (
	TypeSong     Type = "song"
	TypeAlbum    Type = "album"
	TypeArtist   Type = "artist"
	TypePlaylist Type = "playlist"
)

type UrlItem struct {
	Provider string `json:"provider"`
	Type     Type   `json:"type"`
	Id       string `json:"id"`
}
