package model

type PublicConfig struct {
	Instances []struct {
		Url string `json:"url" validate:"required,url"`
	} `json:"instances" validate:"required"`
}

type Config struct {
	Public PublicConfig `json:"public"`
}
