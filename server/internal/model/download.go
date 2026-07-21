package model

import "github.com/google/uuid"

type DownloadStatus string

const (
	DownloadStatusPending DownloadStatus = "pending"
	DownloadStatusRunning DownloadStatus = "running"
	DownloadStatusDone    DownloadStatus = "done"
	DownloadStatusFailed  DownloadStatus = "failed"
	DownloadStatusCancel  DownloadStatus = "cancel"
)

type DownloadTaskInfo struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Data          EnrichedSong   `json:"data"`
	Status        DownloadStatus `json:"status"`
	StatusComment string         `json:"status_comment"`
}

type AddDownloadError struct {
	Type   DataType `json:"type"`
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Reason string   `json:"reason"`
}
