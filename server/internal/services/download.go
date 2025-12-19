package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

type downloadManager struct {
	mu    sync.Mutex
	tasks map[uint]map[uint]*downloadTask
	ids   map[uint]uint
}

type downloadTask struct {
	id             uint
	userId         uint
	api            models.Plugin
	songId         string
	quality        string
	songData       models.SongData
	status         models.Status // depend-on a chan to be change
	retryDownload  chan struct{}
	cancelDownload chan struct{}
	ctx            context.Context
	cancel         context.CancelFunc
}

var DownloadManager = downloadManager{
	mu:    sync.Mutex{},
	tasks: make(map[uint]map[uint]*downloadTask),
	ids:   make(map[uint]uint),
}

func (m *downloadManager) generateId(userId uint) uint {
	id, ok := m.ids[userId]
	if !ok {
		id = 0
	}
	id++
	m.ids[userId] = id

	return id
}

func (m *downloadManager) AddArtist(userId uint, api models.Plugin, artistId string, quality string) {
	artist, err := api.Artist(context.Background(), userId, artistId)
	if err != nil {
		fmt.Printf("downloadManager.AddArtist: %v\n", err)
		return
	}
	for _, album := range artist.Albums {
		m.AddAlbum(userId, api, album.Id, quality)
	}
}

func (m *downloadManager) AddAlbum(userId uint, api models.Plugin, albumId string, quality string) {
	album, err := api.Album(context.Background(), userId, albumId)
	if err != nil {
		fmt.Printf("downloadManager.AddAlbum: %v\n", err)
		return
	}
	for _, song := range album.Songs {
		m.AddSong(userId, api, song.Id, quality)
	}
}

func (m *downloadManager) AddSong(userId uint, api models.Plugin, songId string, quality string) {
	taskId := m.generateId(userId)
	ctx, cancel := context.WithCancel(context.Background())

	newTask := &downloadTask{
		id:             taskId,
		userId:         userId,
		api:            api,
		songId:         songId,
		quality:        quality,
		status:         models.StatusPending,
		retryDownload:  make(chan struct{}),
		cancelDownload: make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}

	m.mu.Lock()
	userTasks, ok := m.tasks[userId]
	if !ok {
		userTasks = make(map[uint]*downloadTask)
		m.tasks[userId] = userTasks
	}
	userTasks[taskId] = newTask
	m.mu.Unlock()

	go m.startMaster(newTask)
}

func saveSong(userId uint, reader io.ReadCloser, extension string, data models.SongData) error {
	defer reader.Close()

	user, err := repository.GetUserByID(userId)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}

	root, err := os.OpenRoot(config.DOWNLOAD_FOLDER)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer root.Close()

	if err := root.Mkdir(user.Username, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	root, err = os.OpenRoot(filepath.Join(config.DOWNLOAD_FOLDER, user.Username))
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer root.Close()

	filename := filepath.Join(data.Artist.Name, data.Album.Title, fmt.Sprintf("%d - %s.%s", data.TrackNumber, data.Title, extension))

	if err := root.Mkdir(data.Artist.Name, 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	if err := root.Mkdir(filepath.Join(data.Artist.Name, data.Album.Title), 0755); err != nil && !os.IsExist(err) {
		return fmt.Errorf("saveSong: %w", err)
	}

	file, err := root.Create(filename)
	if err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		path := file.Name()
		file.Close()
		os.Remove(path)
		return fmt.Errorf("saveSong: %w", err)
	}

	if err := utils.FormatMetadata(filepath.Join(root.Name(), filename), data); err != nil {
		return fmt.Errorf("saveSong: %w", err)
	}
	return nil
}

func (m *downloadManager) startMaster(t *downloadTask) {
	defer close(t.retryDownload)
	defer close(t.cancelDownload)

	ctx, cancel := context.WithCancel(context.Background())
	status := make(chan models.Status)
	data := make(chan models.SongData)

	t.status = models.StatusPending
	go func() {
		status <- models.StatusRunning

		reader, extension, err := t.api.Download(ctx, t.userId, t.songId, t.quality, data)
		if err == nil {
			err = saveSong(t.userId, reader, extension, t.songData)
		}

		if err != nil {
			if errors.Is(err, context.Canceled) {
				status <- models.StatusCancel
				return
			} else {
				status <- models.StatusFailed
				fmt.Println("downloadManager.startMaster: goroutine: ", err)
				return
			}
		}

		status <- models.StatusDone
	}()

	for {
		select {
		case res := <-data:
			t.songData = res
		case res := <-status:
			t.status = res
			if res == models.StatusDone {
				cancel()
				return
			}
		case <-t.retryDownload:
			if t.status == models.StatusCancel || t.status == models.StatusFailed {
				ctx, cancel = context.WithCancel(context.Background())
				status = make(chan models.Status)
				data = make(chan models.SongData)

				t.status = models.StatusPending
				go func() {
					status <- models.StatusRunning

					reader, extension, err := t.api.Download(ctx, t.userId, t.songId, t.quality, data)
					if err == nil {
						err = saveSong(t.userId, reader, extension, t.songData)
					}

					if err != nil {
						if errors.Is(err, context.Canceled) {
							status <- models.StatusCancel
						} else {
							status <- models.StatusFailed
							fmt.Println("downloadManager.startMaster: goroutine: ", err)
						}
					}

					status <- models.StatusDone
				}()
			}
		case <-t.cancelDownload:
			cancel()
		case <-t.ctx.Done():
			cancel()
			return
		}
	}
}

func (m *downloadManager) Retry(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	m.mu.Unlock()
	task.retryDownload <- struct{}{}
	return nil
}

func (m *downloadManager) Cancel(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	m.mu.Unlock()
	task.cancelDownload <- struct{}{}
	return nil
}

func (m *downloadManager) Remove(userId uint, taskId uint) error {
	m.mu.Lock()
	task, ok := m.tasks[userId][taskId]
	if !ok {
		return errors.New("download not found")
	}
	task.cancel()
	delete(m.tasks[userId], taskId)
	m.mu.Unlock()
	return nil
}

func (m *downloadManager) List(userId uint) []models.DownloadData {
	m.mu.Lock()
	tasks := make([]models.DownloadData, 0, len(m.tasks))
	for _, value := range m.tasks[userId] {
		tmp := models.DownloadData{
			Id:     value.id,
			Data:   value.songData,
			Api:    value.api.Name(),
			Status: value.status,
		}
		tasks = append(tasks, tmp)
	}
	m.mu.Unlock()
	return tasks
}
