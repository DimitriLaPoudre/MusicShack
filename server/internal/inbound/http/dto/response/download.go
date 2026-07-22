package response

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type DownloadTask struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	CoverURL      string   `json:"cover_url"`
	Title         string   `json:"title"`
	Artists       []string `json:"artists"`
	Album         string   `json:"album"`
	Status        string   `json:"status"`
	StatusComment string   `json:"status_comment"`
}

func DownloadTaskToResponse(task model.DownloadTaskInfo) DownloadTask {
	artists := []string{}
	for _, artist := range task.Data.Artists {
		artists = append(artists, artist.Name)
	}

	return DownloadTask{
		ID:            task.ID.String(),
		UserID:        task.UserID.String(),
		CoverURL:      task.Data.Album.CoverUrl,
		Title:         task.Data.Title,
		Artists:       artists,
		Album:         task.Data.Album.Title,
		Status:        string(task.Status),
		StatusComment: task.StatusComment,
	}
}

func DownloadTasksToResponse(tasks []model.DownloadTaskInfo) []DownloadTask {
	r := []DownloadTask{}
	for _, task := range tasks {
		r = append(r, DownloadTaskToResponse(task))
	}
	return r
}
