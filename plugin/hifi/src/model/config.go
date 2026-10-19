package model

type PublicConfig struct {
	Instances []struct {
		Url string `json:"url"`
	} `json:"instances"`
}

type Config struct {
	Public PublicConfig `json:"public"`
}
