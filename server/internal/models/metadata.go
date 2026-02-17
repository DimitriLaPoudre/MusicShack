package models

const (
	TagTitle        string = "TITLE"
	TagArtists      string = "ARTISTS"
	TagAlbumArtists string = "ALBUMARTISTS"
	TagAlbum        string = "ALBUM"
	TagTrackNumber  string = "TRACKNUMBER"
	TagVolumeNumber string = "DISCNUMBER"
	TagReleaseDate  string = "RELEASEDATE"
	TagExplicit     string = "ITUNESADVISORY"
	TagAlbumGain    string = "REPLAYGAIN_ALBUM_GAIN"
	TagAlbumPeak    string = "REPLAYGAIN_ALBUM_PEAK"
	TagTrackGain    string = "REPLAYGAIN_TRACK_GAIN"
	TagTrackPeak    string = "REPLAYGAIN_TRACK_PEAK"
	TagISRC         string = "ISRC"
)

func MetadataInfoToTags(info MetadataInfo) map[string][]string {
	return map[string][]string{
		TagTitle:        {info.Title},
		TagAlbum:        {info.Title},
		TagAlbumArtists: info.AlbumArtists,
		TagArtists:      info.Artists,
		TagTrackNumber:  {info.TrackNumber},
		TagVolumeNumber: {info.VolumeNumber},
		TagReleaseDate:  {info.ReleaseDate},
		TagExplicit:     {info.Explicit},
		TagAlbumGain:    {info.AlbumGain},
		TagAlbumPeak:    {info.AlbumPeak},
		TagTrackGain:    {info.TrackGain},
		TagTrackPeak:    {info.TrackPeak},
		TagISRC:         {info.Isrc},
	}
}
