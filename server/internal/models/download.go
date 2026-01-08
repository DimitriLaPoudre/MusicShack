package models

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
	StatusCancel  Status = "cancel"
)

type RequestDownload struct {
	Provider string `json:"provider"`
	Type     string `json:"type"`
	Id       string `json:"id"`
	Quality  string `json:"quality"`
}

type DownloadData struct {
	Id       uint     `json:"id"`
	Provider string   `json:"provider"`
	Data     SongData `json:"data"`
	Status   Status   `json:"status"`
}
