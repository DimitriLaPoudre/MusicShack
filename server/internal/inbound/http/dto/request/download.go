package request

type CreateDownload struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
}
