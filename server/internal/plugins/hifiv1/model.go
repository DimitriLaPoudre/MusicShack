package hifiv1

type HifiV1 struct{}

func (p *HifiV1) Name() string {
	return "hifiV1"
}

func (p *HifiV1) Download(id string, quality string) error {
	return nil
}

func (p *HifiV1) Artist(id string) (any, error) {
	return nil, nil
}
func (p *HifiV1) Search(song, album, artist string) (any, error) {
	return nil, nil
}
func (p *HifiV1) Cover(id string) (string, error) {
	return "", nil
}
func (p *HifiV1) Lyrics(id string) (string, string, error) {
	return "", "", nil
}
