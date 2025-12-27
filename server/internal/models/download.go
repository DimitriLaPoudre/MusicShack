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
	Api     string `json:"api"`
	Type    string `json:"type"`
	Id      string `json:"id"`
	Quality string `json:"quality"`
}

type DownloadData struct {
	Id     uint     `json:"id"`
	Data   SongData `json:"data"`
	Api    string   `json:"api"`
	Status Status   `json:"status"`
}
