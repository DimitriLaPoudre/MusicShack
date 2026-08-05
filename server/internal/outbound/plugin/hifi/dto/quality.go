package dto

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
