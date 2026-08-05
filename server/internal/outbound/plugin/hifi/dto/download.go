package dto

type DownloadItem struct {
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

type ManifestTidal struct {
	MimeType       string   `json:"mimeType"`
	Codecs         string   `json:"codecs"`
	EncryptionType string   `json:"encryptionType"`
	URLs           []string `json:"urls"`
}

type ManifestMPD struct {
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
